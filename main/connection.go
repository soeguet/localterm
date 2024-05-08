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

type PayloadType int

const (
	AuthenticationTypeConst  PayloadType = 0
	MessageTypeConst         PayloadType = 1
	ClientListTypeConst      PayloadType = 2
	MessageListTypeConst     PayloadType = 4
	TypingIndicatorTypeConst PayloadType = 5
	ClientTypingConst        PayloadType = 6
	GetMessagesTypeConst     PayloadType = 7
)

func handlePayloadsOfMessageType(message []byte, app *App) {

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

func handleDesktopNotificationPossibility(messagePayload MessagePayload, app *App) {
	// check if message is from this client and do not send desktop notification
	if messagePayload.ClientType.ClientDbId == envVars.Id {
		return
	}

	// send desktop notification
	if err := app.desktopNotification(&messagePayload); err != nil {
		fmt.Println("Error sending desktop notification for received message:", err)
		return
	}
}

func unmarshallPayloadToMessagePayload(message []byte) (MessagePayload, error) {
	var messagePayload MessagePayload

	if err := json.Unmarshal(message, &messagePayload); err != nil {
		fmt.Println("Error parsing messagePayload:", err)
		return messagePayload, err
	}
	return messagePayload, nil
}

func handlePayloadsOfMessageListType(message []byte, app *App) {

	messageListPayload := unmarshallMessageToMessageListPayload(message)

	resetMessageCache()
	app.clearChatView()

	for _, payload := range messageListPayload.MessageList {
		index := appendMessageToCache(payload)
		addNewMessageToScrollPanel(&index, &payload)
	}
}

func unmarshallMessageToMessageListPayload(message []byte) MessageListPayload {
	var messageListPayload MessageListPayload
	if err := json.Unmarshal(message, &messageListPayload); err != nil {
		fmt.Println("Error parsing messageListPayload:", err)
	}
	return messageListPayload
}

func handlePayloadsOfTypingIndicatorType(message []byte, app *App) {
	var typingPayload TypingPayload
	if err := json.Unmarshal(message, &typingPayload); err != nil {
		fmt.Println("Error parsing typingPayload:", err)
	}

	// TODO implement typing indicator
}

func handlePayloadsOfClientListType(message []byte, app *App) {

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

func unmarshallMessageToClientListPayload(message []byte) (ClientList, error) {
	var clientListPayload ClientList
	if err := json.Unmarshal(message, &clientListPayload); err != nil {
		fmt.Println("Error parsing clientListPayload:", err)
		return clientListPayload, err
	}
	return clientListPayload, nil
}

func handlePayload(message []byte, app *App) {
	var msg GenericMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	switch msg.PayloadType {

	case MessageTypeConst:
		handlePayloadsOfMessageType(message, app)

	case ClientListTypeConst:
		handlePayloadsOfClientListType(message, app)

	case MessageListTypeConst:
		handlePayloadsOfMessageListType(message, app)

	case TypingIndicatorTypeConst:
		handlePayloadsOfTypingIndicatorType(message, app)

	case GetMessagesTypeConst:
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
	messagePayload := MessagePayload{
		PayloadType: MessageTypeConst,
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

	// Send the message
	err := conn.WriteJSON(messagePayload)
	if err != nil {
		fmt.Println("Error writing messagePayload:", err)
	}
}

func retrieveLast100Messages(c *websocket.Conn) {
	// Get the last 100 messages
	messageListPayload := MessageListRequestPayload{
		PayloadType: MessageListTypeConst,
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
	authenticationPayload := AuthenticationPayload{
		PayloadType:    AuthenticationTypeConst,
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
