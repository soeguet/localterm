// main package
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

type payloadType int

const (
	authenticationTypeConst  payloadType = 0
	messageTypeConst         payloadType = 1
	clientListTypeConst      payloadType = 2
	messageListTypeConst     payloadType = 4
	typingIndicatorTypeConst payloadType = 5
	clientTypingConst        payloadType = 6
	getMessagesTypeConst     payloadType = 7
)

func handlePayloadsOfMessageType(message []byte, app *app) {
	messagePayload, err := unmarshallPayloadToMessagePayload(message)
	if err != nil {
		fmt.Println("Error parsing messagePayload:", err)
		if err := app.notifier.Notify("Error", "Error parsing messagePayload", ""); err != nil {
			log.Fatalf("Error sending desktop notification for error parsing messagePayload: %v", err)
		}
		return
	}

	index := appendMessageToCache(messagePayload)
	addNewMessageToScrollPanel(&index, &messagePayload)

	handleDesktopNotificationPossibility(messagePayload, app)
}

func handleDesktopNotificationPossibility(messagePayload messagePayload, app *app) {
	// check if message is from this client and do not send desktop notification
	if messagePayload.ClientType.ClientDbID == envVars.ID {
		return
	}

	// send desktop notification
	if err := app.desktopNotification(&messagePayload); err != nil {
		fmt.Println("Error sending desktop notification for received message:", err)
		return
	}
}

func unmarshallPayloadToMessagePayload(message []byte) (messagePayload, error) {
	var messagePayload messagePayload

	if err := json.Unmarshal(message, &messagePayload); err != nil {
		fmt.Println("Error parsing messagePayload:", err)
		return messagePayload, err
	}
	return messagePayload, nil
}

func handlePayloadsOfMessageListType(message []byte, app *app) {
	messageListPayload := unmarshallMessageToMessageListPayload(message)

	resetMessageCache()
	app.clearChatView()

	for _, payload := range messageListPayload.MessageList {
		index := appendMessageToCache(payload)
		addNewMessageToScrollPanel(&index, &payload)
	}
}

func unmarshallMessageToMessageListPayload(message []byte) messageListPayload {
	var messageListPayload messageListPayload
	if err := json.Unmarshal(message, &messageListPayload); err != nil {
		fmt.Println("Error parsing messageListPayload:", err)
	}
	return messageListPayload
}

func handlePayloadsOfTypingIndicatorType(message []byte, app *app) {
	var typingPayload typingPayload
	if err := json.Unmarshal(message, &typingPayload); err != nil {
		fmt.Println("Error parsing typingPayload:", err)
	}

	if typingPayload.IsTyping {
		addTypingClient(typingPayload.ClientDbID)
	} else {
		removeTypingClient(typingPayload.ClientDbID)
	}

	typingLabelText := generateTypingString()
	fmt.Println(typingLabelText)
	app.setTypingLabelText(typingLabelText)
}

func handlePayloadsOfClientListType(message []byte, app *app) {
	clientListPayload, err := unmarshallMessageToClientListPayload(message)
	if err != nil {
		fmt.Println("Error parsing clientListPayload:", err)
		if err := app.notifier.Notify("Error", "Error parsing clientListPayload", ""); err != nil {
			log.Fatalf("Error sending desktop notification for error parsing clientListPayload: %v", err)
		}
		return
	}

	setClientList(&clientListPayload)
	retrieveLast100Messages(app.conn)
}

func unmarshallMessageToClientListPayload(message []byte) (clientListStruct, error) {
	var clientListPayload clientListStruct
	if err := json.Unmarshal(message, &clientListPayload); err != nil {
		fmt.Println("Error parsing clientListPayload:", err)
		return clientListPayload, err
	}
	return clientListPayload, nil
}

func handlePayload(message []byte, app *app) {
	var msg genericMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	switch msg.PayloadType {

	case messageTypeConst:
		handlePayloadsOfMessageType(message, app)

	case clientListTypeConst:
		handlePayloadsOfClientListType(message, app)

	case messageListTypeConst:
		handlePayloadsOfMessageListType(message, app)

	case typingIndicatorTypeConst:
		handlePayloadsOfTypingIndicatorType(message, app)

	case getMessagesTypeConst:
		retrieveLast100Messages(app.conn)

	default:
		fmt.Println("unknown PayloadType", msg.PayloadType)
	}

	app.ui.Draw()
}

func connection(app *app) error {
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
	messagePayload := messagePayload{
		PayloadType: messageTypeConst,
		MessageType: messageType{
			MessageDbID:    GenerateRandomID(),
			MessageContext: base64.StdEncoding.EncodeToString([]byte(*message)),
			Deleted:        false,
			Edited:         false,
			MessageTime:    time.Now().Format("15:04"),
			MessageDate:    time.Now().Format("2006-01-02"),
		},
		ClientType: clientType{
			ClientDbID: envVars.ID,
		},
	}

	// Send the message
	err := conn.WriteJSON(messagePayload)
	if err != nil {
		fmt.Println("Error writing messagePayload:", err)
	}
}

func retrieveLast100Messages(c *websocket.Conn) {
	// Get the last 100 messages
	messageListPayload := messageListRequestPayload{
		PayloadType: messageListTypeConst,
	}

	if err := c.WriteJSON(messageListPayload); err != nil {
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
	authenticationPayload := authenticationPayload{
		PayloadType:    authenticationTypeConst,
		ClientUsername: getEnvUsername(),
		ClientDbID:     envVars.ID,
	}
	return json.Marshal(authenticationPayload)
}

func createConnection(localChatIP string, localChatPort string) (*websocket.Conn, error) {
	url := fmt.Sprintf("ws://%s:%s/chat", localChatIP, localChatPort)

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
