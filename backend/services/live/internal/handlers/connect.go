package handlers

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

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// func (h *Handlers) Connect() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		ctx := r.Context()

// 		fmt.Println("start websocket handler called")

// 		websocketConn, cancelConn, err := upgradeConnection(w, r)
// 		if err != nil {
// 			tele.Error(ctx, "failed to upgrade connection to websocket")
// 			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "something went wrong. Error: E314759")
// 			return
// 		}
// 		defer cancelConn()

// 		claims, ok := utils.GetValue[jwt.Claims](r, ct.ClaimsKey)
// 		if !ok {
// 			tele.Error(ctx, "failed to fetch claims")
// 			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "failed to fetch claims")
// 		}

// 		clientId := claims.UserId
// 		fmt.Println("Client id:", clientId)

// 		// context for cancelling send goroutine when connection closes, and for other goroutines to not send anything to channel that is no longer needed
// 		ctx, cancelContext := context.WithCancel(context.Background())

// 		connectionId, ok := ctx.Value(ct.ReqID).(string)
// 		if !ok {
// 			tele.Error(ctx, "failed to fetch claims")
// 			utils.ErrorJSON(ctx, w, http.StatusInternalServerError, "Something went wrong. Error: E741786")

// 		}
// 		channel := make(chan (WSMessage))

// 		// wait group to wait for sender goroutine to stop
// 		var wg sync.WaitGroup

// 		wg.Add(1)
// 		go websocketSender(connectionId, clientId, channel, websocketConn, ctx, &wg)

// 		websocketListener(websocketConn, clientId, connectionId)

// 		// CLOSING PRECEDURE

// 		// remove client so that no other goroutine attempts to send here
// 		// clients.RemoveClientConnection(clientId, connectionId)

// 		//adding message about online status here, may need to move
// 		// sendUserStatusUpdate(clientId, "offline")

// 		// remove all subscriptions associated with connectionId
// 		// eventSystem.PurgeSubscriptions(connectionId)

// 		// close the context to stop sender goroutine and other goroutines from sending
// 		cancelContext()

// 		// wait for sender goroutine to stop
// 		wg.Wait()

// 		fmt.Println("ws handler closing")
// 	}
// }

// func sendErrorToWS(websocketConn *websocket.Conn, payload string) {
// 	errorMessage := []any{WSMessage{
// 		Payload: payload,
// 	}}

// 	bundledMessage, err := json.Marshal(errorMessage)
// 	if err != nil {
// 		fmt.Println("this isn't supposed to happen")
// 		panic(1)
// 	}

// 	err = websocketConn.WriteMessage(websocket.TextMessage, bundledMessage)
// 	if err != nil {
// 		fmt.Println("failed to inform user that they have too many tabs open")
// 	}
// }

// // routine that reads data coming from this client connection
// func websocketListener(websocketConn *websocket.Conn, clientId int64, connectionId int64) {
// 	for {
// 		_, msg, err := websocketConn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
// 				break
// 			}
// 			log.Printf("clientId: %d, Start websocket error: unexpected read error: %v\n", clientId, err)
// 			return
// 		}

// 		messageString := string(msg)
// 		if len(messageString) < 2 {
// 			fmt.Println("invalid message received:", messageString)
// 			return
// 		}

// 	}
// }

// // Goroutine that sends data to this connection
// func websocketSender(connectionid string, clientId int64, channel <-chan WSMessage, conn *websocket.Conn, ctx context.Context, wg *sync.WaitGroup) {
// 	var timer *time.Timer
// 	defer wg.Done()
// 	for {
// 		select {
// 		case message := <-channel:
// 			slice := []WSMessage{message}
// 			bundledMessage, err := json.Marshal(slice)
// 			if err != nil {
// 				fmt.Println("AAAAA! A BAD JSON WAS GIVEN! The double encode failed! In production this situation will probably be ignored")
// 				panic(1)
// 				// TODO attempt individual sends? and drop w/e fails with a warning message
// 			}

// 			fmt.Println("sending", string(bundledMessage), "to", connectionid)
// 			err = conn.WriteMessage(websocket.TextMessage, bundledMessage)
// 			if err != nil {
// 				fmt.Printf("connectionid: %d, error on send message: clientid:%d err:%s \n", connectionid, clientId, err.Error())
// 				break outerLoop
// 			}
// 			break bundleLoop

// 		case <-ctx.Done():
// 			if timer != nil {
// 				if !timer.Stop() {
// 					<-timer.C
// 				}
// 			}
// 			break outerLoop
// 		}
// 	}
// 	// TODO look into errors that can only happen when writing, this should cancel the entire connection
// }

// func upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, func(), error) {
// 	websocketConn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		utils.ErrorJSON(r.Context(), w, http.StatusInternalServerError, "failed websocket upgrade")
// 		return nil, func() {}, fmt.Errorf("failed to upgrade connection to websocket: %w", err)
// 	}
// 	deferMe := func() {
// 		err := websocketConn.Close()
// 		if err != nil {
// 			fmt.Println("failed to close websocket connection")
// 		}
// 	}
// 	return websocketConn, deferMe, nil
// }

// type WSMessage struct {
// 	Payload any `json:"message"`
// }
