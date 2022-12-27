package handlers

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin:     func(r *http.Request) bool { return true },
// }

// func ConnectHandler(w http.ResponseWriter, r *http.Request) {
// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer ws.Close()

// 	device := waContainer.NewDevice()
// 	waClient := whatsmeow.NewClient(device, nil)

// 	qrCh, _ := waClient.GetQRChannel(r.Context())
// 	if err := waClient.Connect(); err != nil {
// 		log.Fatal(err)
// 	}

// 	for item := range qrCh {
// 		if item.Event == "code" {
// 			fmt.Printf("%v", item.Timeout)
// 			ws.WriteMessage(websocket.TextMessage, []byte(item.Code))
// 		} else {
// 			ws.WriteMessage(websocket.TextMessage, []byte("connected"))
// 			break
// 		}
// 	}
// }
