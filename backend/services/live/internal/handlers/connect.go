package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"social-network/shared/gen-go/chat"
	"social-network/shared/go/batching"
	"social-network/shared/go/ct"
	utils "social-network/shared/go/http-utils"
	"social-network/shared/go/jwt"
	tele "social-network/shared/go/telemetry"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

//TODO add redis call to ratelimit the number of open connections

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
		ctx := r.Context()
		tele.Info(ctx, "start websocket handler called")
		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
		if !ok {
			tele.Error(ctx, "failed to fetch claims")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to fetch claims")
			return
		}
		clientId := claims.UserId

		connectionId, ok := ctx.Value(ct.ReqID).(string)
		if !ok {
			tele.Error(ctx, "failed get request id")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Something went wrong. Error: E741786")
			return
		}

		// UPGRADE
		websocketConn, cancelConn, err := upgradeConnection(ctx, w, r)
		if err != nil {
			tele.Error(ctx, "failed to upgrade connection to websocket")
			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Something went wrong. Error: E314759")
			return
		}
		defer cancelConn()

		wsChannel := make(chan []byte) //send listens to this and forwards messages to user using the websocket

		natsHandler := func(msg *nats.Msg) { //handles nats messages, just forwards them to the above channel
			wsChannel <- msg.Data
			tele.Info(ctx, "forwarded nats message to websocket @1", "connection", connectionId)
		}

		var wg sync.WaitGroup
		wg.Go(func() { h.websocketSender(ctx, wsChannel, websocketConn) })
		h.websocketListener(ctx, websocketConn, clientId, connectionId, natsHandler)

		wg.Wait()

		tele.Info(ctx, "ws handler closing")
	}
}

// routine that reads data coming from this client connection, reads the message and decides what to do with it
func (h *Handlers) websocketListener(ctx context.Context, websocketConn *websocket.Conn, clientId int64, connectionId string, handler nats.MsgHandler) {
	subcriptions := make(map[string]*nats.Subscription)
	tele.Info(ctx, "websocket listener started for connection @1", "connection", connectionId)
	key := ct.PrivateMessageKey(clientId)
	sub, err := h.Nats.Subscribe(key, handler)
	tele.Info(ctx, "subscribed to conversation @1 using key @2", "conversation", clientId, "key", key)
	if err != nil {
		tele.Error(ctx, "websocket subscription @1", "error", err.Error())
		return
	}
	subcriptions[key] = sub

	defer func() { //unsub from all
		for _, sub := range subcriptions {
			err := sub.Unsubscribe()
			if err != nil {
				tele.Error(ctx, "websocket unsubscribe @1", "error", err.Error())
			}
		}
	}()

	for {
		_, msg, err := websocketConn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				tele.Info(ctx, "websocket closed")
				break
			}
			tele.Error(ctx, "Start websocket error, unexpected read:  @1", "error", err.Error())
			return
		}

		messageString := string(msg)
		if len(messageString) < 2 {
			tele.Error(ctx, "invalid message received: @1 from @2", "message", messageString, "connection", connectionId)
			continue
		}

		tele.Info(ctx, "received: @1 from @2", "message", messageString, "connection", connectionId)

		parts := strings.SplitN(messageString, ":", 2)
		if len(parts) != 2 {
			tele.Error(ctx, "malformed message received: @1 from @2, @3", "message", messageString, "connection", connectionId, "parts", parts)
			continue
		}

		msgType := parts[0]
		payload := parts[1]

		switch msgType {
		case "conversation":
			//TODO verify that they are allowed to subscribe there
			tele.Info(ctx, "subscribing to conversation @1", "conversation", payload)
			sub, err := h.Nats.Subscribe(payload, handler)
			if err != nil {
				tele.Error(ctx, "websocket subscription @1", "error", err.Error())
				continue
			}
			subcriptions[payload] = sub
		case "unsub":
			sub := subcriptions[payload]
			tele.Info(ctx, "unsubscribing from conversation @1", "conversation", payload)
			err := sub.Unsubscribe()
			if err != nil {
				tele.Error(ctx, "websocket unsubscribe @1", "error", err.Error())
				continue
			}
			delete(subcriptions, payload)
		case "ch":

			_, err = h.ChatService.GetOrCreatePrivateConv(ctx, &chat.GetOrCreatePrivateConvRequest{
				User:              clientId,
				OtherUser:         2,
				RetrieveOtherUser: false,
			})
			if err != nil {
				tele.Error(ctx, "failed to get or create private conversation @1", "error", err.Error())
			}

			type chatMessage struct {
				Category       string     `json:"category"`
				ConversationId ct.Id      `json:"conversation_id"`
				Body           ct.MsgBody `json:"body"`
			}
			message := &chatMessage{}
			err = json.Unmarshal([]byte(payload), message)
			if err != nil {
				tele.Error(ctx, "failed to unmarshal chat message @1", "error", err.Error())
			}
			_, err = h.ChatService.CreatePrivateMessage(ctx, &chat.CreatePrivateMessageRequest{
				ConversationId: message.ConversationId.Int64(),
				SenderId:       clientId,
				MessageText:    message.Body.String(),
			})
			if err != nil {
				tele.Error(ctx, "failed to create private message @1", "error", err.Error())
			}
		}
	}
}

// Goroutine that sends data to this connection, it can pool messages if they arrive fast enough
func (h *Handlers) websocketSender(ctx context.Context, channel <-chan []byte, conn *websocket.Conn) {
	payloadBytes := []byte{}

	//handler is given to batcher, so that the batcher calls it with many accumulated messages at once
	handler := func(messages []json.RawMessage) error {
		var err error
		tele.Info(ctx, "about to marshal and send @1 messages to websocket, @2", "count", len(messages), "rawBody", messages)
		payloadBytes, err = json.Marshal(messages)
		if err != nil {
			return err
		}
		err = conn.WriteMessage(websocket.TextMessage, payloadBytes)
		if err != nil {
			return err
		}
		return nil
	}

	batcherInput, errChannel := batching.Batcher(ctx, handler, time.Millisecond*200, 200)

	for {
		select {
		case message, ok := <-channel:
			if !ok {
				tele.Error(ctx, "websocket channel closed")
				return //TODO check if ok
			}
			batcherInput <- json.RawMessage(message)

		case err := <-errChannel:
			tele.Error(ctx, "(batched) write message @1", "error", err.Error())
			//TODO figure out if something needs to be done here

		case <-ctx.Done():
			tele.Info(ctx, "websocket context ended")
			return
		}
	}
}

// func sendErrorToWS(ctx context.Context, websocketConn *websocket.Conn, payload string) {
// 	errorMessage := []string{payload}

// 	bundledMessage, err := json.Marshal(errorMessage)
// 	if err != nil {
// 		tele.Error(ctx, "this isn't supposed to happen")
// 		panic(1) //???
// 	}

// 	err = websocketConn.WriteMessage(websocket.TextMessage, bundledMessage)
// 	if err != nil {
// 		tele.Info(ctx, "failed to inform user that they have too many tabs open")
// 	}
// }

func upgradeConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) (*websocket.Conn, func(), error) {
	websocketConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.ErrorJSON(r.Context(), w, http.StatusInternalServerError, "failed websocket upgrade")
		tele.Warn(ctx, "failed to upgrade websocket @1", "error", err.Error())
		return nil, func() {}, fmt.Errorf("failed to upgrade connection to websocket: %w", err)
	}
	deferMe := func() {
		err := websocketConn.Close()
		if err != nil {
			tele.Warn(ctx, "failed to close websocket connection")
		}
	}
	return websocketConn, deferMe, nil
}
