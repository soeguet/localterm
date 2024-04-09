package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/rivo/tview"
)

func Connection(app *tview.Application) error {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:5588/chat", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			payload := DecodeJson(message)

			switch v := payload.(type) {

			case MessagePayload:
				AddNewPlainMessageToChatView("Payload: " + v.MessageType.MessageContext)
			case ClientList:
				AddNewPlainMessageToChatView("ClientList:")
			case MessageListPayload:
				for i, _ := range v.MessageList {

					AddNewPlainMessageToChatView(string("henlos" + string(i)))

				}
			case int:
				if v == -1 {
					fmt.Println("unknown type")
				}
			default:
				fmt.Println("unknows type")
			}

		}
	}()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Println("interrupt")

			// Clean disconnect
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return nil
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}