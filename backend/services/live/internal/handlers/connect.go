package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"social-network/shared/go/ct"
	utils "social-network/shared/go/http-utils"
	"social-network/shared/go/jwt"
	tele "social-network/shared/go/telemetry"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// FOUND FROM DOCUMENTATION:

// The message types are defined in RFC 6455, section 11.8.

// TextMessage denotes a text data message. The text message payload is
// interpreted as UTF-8 encoded text data.
// TextMessage = 1

// BinaryMessage denotes a binary data message.
// BinaryMessage = 2

// CloseMessage denotes a close control message. The optional message
// payload contains a numeric code and text. Use the FormatCloseMessage
// function to format a close message payload.
// CloseMessage = 8

// PingMessage denotes a ping control message. The optional message payload
// is UTF-8 encoded text.
// PingMessage = 9

// PongMessage denotes a pong control message. The optional message payload
// is UTF-8 encoded text.
// PongMessage = 10

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (h *Handlers) Connect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("start websocket handler called")

		ctx := r.Context()
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "failed to fetch claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to fetch claims")
		}
		clientId := claims.UserId

		// UPGRADE
		websocketConn, cancelConn, err := upgradeConnection(w, r)
		if err != nil {
			tele.Error(ctx, "failed to upgrade connection to websocket")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Something went wrong. Error: E314759")
			return
		}
		defer cancelConn()

		connectionId, ok := ctx.Value(ct.ReqID).(string)
		if !ok {
			tele.Error(ctx, "failed get request id")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Something went wrong. Error: E741786")
		}

		channel := make(chan (string))

		var wg sync.WaitGroup
		wg.Go(func() { websocketSender(ctx, connectionId, clientId, channel, websocketConn) })
		websocketListener(websocketConn, clientId, connectionId)
		wg.Wait()

		fmt.Println("ws handler closing")
	}
}

// routine that reads data coming from this client connection
func websocketListener(websocketConn *websocket.Conn, clientId int64, connectionId string) {
	for {
		_, msg, err := websocketConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				break
			}
			log.Printf("clientId: %d, Start websocket error: unexpected read error: %v\n", clientId, err)
			return
		}

		messageString := string(msg)
		if len(messageString) < 2 {
			fmt.Println("invalid message received:", messageString)
			return
		}

		//do something with message received
		//forward to chat service
		//send to chat client
	}
}

// Goroutine that sends data to this connection, it can pool messages if they arrive fast enough
func websocketSender(ctx context.Context, connectionid string, clientId int64, channel <-chan string, conn *websocket.Conn) {
	timer := time.NewTimer(time.Hour)
	timer.Stop()
	timerOn := false
	lastFlushTimeStamp := time.Now().Add(-time.Hour)
	poolingDuration := time.Millisecond * 500

	messageBucket := []json.RawMessage{}
	payloadBytes := []byte{}
	var err error
	for {
		select {
		case message := <-channel:
			// message arrived, adding it to bucket
			messageBucket = append(messageBucket, json.RawMessage(message))

			if timerOn {
				// timer is already activated, therefore we just gather messages until the timer fires
				continue
			}

			//check if more than 'poolingDuration' amount of time has passed after the last flush
			if time.Since(lastFlushTimeStamp) <= poolingDuration {
				//new message came too soon, so we'll just start the timer and wait it to ring before we flush
				timer.Reset(poolingDuration - time.Since(lastFlushTimeStamp))
				timerOn = true
				continue
			}

			// timer is off, and it's been a while since the last message, lets sending it immediately

			// preparing message
			payloadBytes, err = json.Marshal(messageBucket)
			if err != nil {
				fmt.Println("ERROR")
				return
			}

			//sending message
			err = conn.WriteMessage(websocket.TextMessage, payloadBytes)
			if err != nil {
				fmt.Printf("connectionid: %s, error on send message: clientid:%d err:%s \n", connectionid, clientId, err.Error())
				return
			}

			//clear the bucket
			messageBucket = []json.RawMessage{}
			lastFlushTimeStamp = time.Now()

			// activating timer so that if another message comes too soon,
			// we gather it into the bucket instead of sending it immediately
			timer.Reset(poolingDuration)
			timerOn = true

		case <-timer.C:
			timerOn = false
			if len(messageBucket) == 0 {
				// bucket is empty, therefore no need to send anything
				continue
			}

			//preparing message
			payloadBytes, err = json.Marshal(messageBucket)
			if err != nil {
				fmt.Println("ERROR")
				return
			}

			// bucket has piled up messages, so it's time to flush
			err = conn.WriteMessage(websocket.TextMessage, payloadBytes)
			if err != nil {
				fmt.Printf("connectionid: %s, error on send message: clientid:%d err:%s \n", connectionid, clientId, err.Error())
				return
			}

			//clearing bucket
			messageBucket = []json.RawMessage{}
			lastFlushTimeStamp = time.Now()

		case <-ctx.Done():
			timer.Stop()
			return
		}
	}
}

func sendErrorToWS(websocketConn *websocket.Conn, payload string) {
	errorMessage := []string{}

	bundledMessage, err := json.Marshal(errorMessage)
	if err != nil {
		fmt.Println("this isn't supposed to happen")
		panic(1)
	}

	err = websocketConn.WriteMessage(websocket.TextMessage, bundledMessage)
	if err != nil {
		fmt.Println("failed to inform user that they have too many tabs open")
	}
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, func(), error) {
	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.ErrorJSON(r.Context(), w, http.StatusInternalServerError, "failed websocket upgrade")
		return nil, func() {}, fmt.Errorf("failed to upgrade connection to websocket: %w", err)
	}
	deferMe := func() {
		err := websocketConn.Close()
		if err != nil {
			fmt.Println("failed to close websocket connection")
		}
	}
	return websocketConn, deferMe, nil
}
