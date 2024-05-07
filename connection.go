package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

func handlePayload(message []byte, app *App) {
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
		index := appendMessageToCache(messagePayload)
		addNewMessageViaMessagePayload(&index, &messagePayload)
	case 2:
		var clientListPayload ClientList
		if err := json.Unmarshal(message, &clientListPayload); err != nil {
			fmt.Println("Error parsing clientListPayload:", err)
		}
		setClientList(&clientListPayload)
		retrieveLast100Messages(app.conn)
	case 4:
		var messageListPayload MessageListPayload
		if err := json.Unmarshal(message, &messageListPayload); err != nil {
			fmt.Println("Error parsing messageListPayload:", err)
		}
		resetMessageCache()
		app.clearChatView()
		for _, payload := range messageListPayload.MessageList {
			index := appendMessageToCache(payload)
			addNewMessageViaMessagePayload(&index, &payload)
		}
	case 5:
		var typingPayload TypingPayload
		if err := json.Unmarshal(message, &typingPayload); err != nil {
			fmt.Println("Error parsing typingPayload:", err)
		}
	case 7:
		retrieveLast100Messages(app.conn)
	default:
		fmt.Println("unknown PayloadType", msg.PayloadType)
	}
	app.ui.Draw()
}

func connection(app *App) error {
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

			handlePayload(message, app)
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
			// Wait for the server to close the connection after sending the close message
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
			MessageDbId:    "TOBEREMOVED",
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
}

func authenticateClientAtSocket(c *websocket.Conn) error {
	authenticationPayloadBytes, err := getAuthenticationPayloadBytes()
	if err != nil {
		return fmt.Errorf("error marshalling authenticationPayload: %v", err)
	}

	writeError := func(err error) error {
		if err != nil {
			return fmt.Errorf("error writing authenticationPayload: %v", err)
		}
		return nil
	}
	return writeError(c.WriteMessage(websocket.TextMessage, authenticationPayloadBytes))
}

func getAuthenticationPayloadBytes() ([]byte, error) {
	authenticationPayload := AuthenticationPayload{
		PayloadType:    0,
		ClientUsername: getEnvUsername(),
		ClientDbId:     envVars.Id,
	}
	return json.Marshal(authenticationPayload)
}

func createConnection(localChatIp string, localChatPort string) (*websocket.Conn, error) {
	url := fmt.Sprintf("ws://%s:%s/chat", localChatIp, localChatPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
