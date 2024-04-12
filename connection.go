package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateBase36(n int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, n)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func generateSimpleId() string {
	timestamp := time.Now().UnixMilli()
	randomPart := generateBase36(7) // 7 Zeichen, um der Länge in Deinem TypeScript-Beispiel zu entsprechen
	return fmt.Sprintf("id-%d-%s", timestamp, randomPart)
}

func Connection(app *App) error {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	defer func(c *websocket.Conn) {
		err := c.Close()
		if err != nil {
			log.Fatal("close:", err)
		}
	}(app.conn)

	done := make(chan struct{})

	go func() {
		defer close(done)

		for {

			_, message, err := app.conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			var msg GenericMessage
			if err := json.Unmarshal(message, &msg); err != nil {
				fmt.Println("Error parsing JSON:", err)
				return
			}

			switch msg.PayloadType {

			case 1:
				var messagePayload MessagePayload
				if err := json.Unmarshal(message, &messagePayload); err != nil {
					fmt.Println("Error parsing messagePayload:", err)
				}
				// AddNewEncryptedMessageToChatView(&messagePayload.MessageType.MessageContext)
				AddNewMessageViaMessagePayload(&messagePayload)

			case 2:
				var clientListPayload ClientList
				if err := json.Unmarshal(message, &clientListPayload); err != nil {
					fmt.Println("Error parsing clientListPayload:", err)
				}
				SetClientList(&clientListPayload)

				// request the last 100 messages AFTER the client list is received, otherwise race condition // TODO fix this
				retrieveLast100Messages(app.conn)

			case 4:
				var messageListPayload MessageListPayload

				if err := json.Unmarshal(message, &messageListPayload); err != nil {
					fmt.Println("Error parsing messageListPayload:", err)
				}
				for _, payload := range messageListPayload.MessageList {
					AddNewMessageViaMessagePayload(&payload)
				}

			case 5:
				var typingPayload TypingPayload
				if err := json.Unmarshal(message, &typingPayload); err != nil {
					fmt.Println("Error parsing typingPayload:", err)
				}
				AddNewPlainMessageToChatView(&typingPayload.ClientDbId)

			default:
				fmt.Println("unknown PayloadType", msg.PayloadType)
			}

			app.ui.Draw()
		}
	}()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			log.Println("interrupt")

			// Clean disconnect
			err := app.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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

func sendMessagePayloadToWebsocket(conn *websocket.Conn, message *string) {
	// Send the message
	messagePayload := MessagePayload{
		PayloadType: 1,
		MessageType: MessageType{
			MessageDbId:    generateSimpleId(),
			MessageContext: base64.StdEncoding.EncodeToString([]byte(*message)),
			MessageTime:    time.Now().Format("15:04"),
			MessageDate:    time.Now().Format("2006-01-02"),
		},
		ClientType: ClientType{
			ClientDbId: envVars.Id,
		},
	}

	err := conn.WriteJSON(messagePayload)
	if err != nil {
		fmt.Println("Error writing messagePayload:", err)
	}
}

func retrieveLast100Messages(c *websocket.Conn) {
	// Get the last 100 messages
	messageListPayload := MessageListRequestPayload{
		PayloadType: 4,
	}
	err := c.WriteJSON(messageListPayload)
	if err != nil {
		fmt.Println("Error writing messageListPayload:", err)
	}

	// messageListPayloadBytes, err := json.Marshal(messageListPayload)
	// if err != nil {
	// 	fmt.Println("Error marshalling messageListPayload:", err)
	// }
	// err = c.WriteMessage(websocket.TextMessage, messageListPayloadBytes)
	// if err != nil {
	// 	fmt.Println("Error writing messageListPayload:", err)
	// }
}

func authenticateClientAtSocket(c *websocket.Conn) {
	// Authenticate the client
	authenticationPayload := AuthenticationPayload{
		PayloadType:    0,
		ClientUsername: envVars.Username,
		ClientDbId:     envVars.Id,
	}
	authenticationPayloadBytes, err := json.Marshal(authenticationPayload)
	if err != nil {
		fmt.Println("Error marshalling authenticationPayload:", err)
	}
	err = c.WriteMessage(websocket.TextMessage, authenticationPayloadBytes)
	if err != nil {
		fmt.Println("Error writing authenticationPayload:", err)
	}
}

func requestClientList(c *websocket.Conn) {
	// Get the client list
	clientListPayload := ClientListRequestPayload{
		PayloadType: 2,
	}
	clientListPayloadBytes, err := json.Marshal(clientListPayload)
	if err != nil {
		fmt.Println("Error marshalling clientListPayload:", err)
	}
	err = c.WriteMessage(websocket.TextMessage, clientListPayloadBytes)
	if err != nil {
		fmt.Println("Error writing clientListPayload:", err)
	}
}
