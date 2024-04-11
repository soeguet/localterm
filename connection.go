package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"

	"github.com/rivo/tview"
)

/**
 * [[ RESULTING TYPE ]]
 * export type MessagePayload = {
 *      payloadType: PayloadSubType.message;
 *      messageType: {
 *          messageDbId: string;
 *          messageContext: string;
 *          messageTime: string;
 *          messageDate: Date;
 *      };
 *      clientType: {
 *          clientDbId: string;
 *      };
 *      quoteType?: {
 *          quoteMessageId: string;
 *          quoteClientId: string;
 *          quoteMessageContext: string;
 *          quoteTime: string;
 *          quoteDate: Date;
 *      };
 *      reactionType?: {
 *          reactionMessageId: string;
 *          reactionContext: string;
 *          reactionClientId: string;
 *      }[];
 *    };
 */
type GenericMessage struct {
	PayloadType int             `json:"payloadType"`
	Payload     json.RawMessage `json:"payload"`
}

type MessagePayload struct {
	PayloadType  int            `json:"payloadType"`
	MessageType  MessageType    `json:"messageType"`
	ClientType   ClientType     `json:"clientType"`
	QuoteType    QuoteType      `json:"quoteType"`
	ReactionType []ReactionType `json:"reactionType"`
}

type MessageType struct {
	MessageDbId    string `json:"messageDbId"`
	MessageContext string `json:"messageContext"`
	MessageTime    string `json:"messageTime"`
	MessageDate    string `json:"messageDate"`
}

type ClientType struct {
	ClientDbId string `json:"clientDbId"`
}

type QuoteType struct {
	QuoteMessageId      string `json:"quoteMessageId"`
	QuoteClientId       string `json:"quoteClientId"`
	QuoteMessageContext string `json:"quoteMessageContext"`
	QuoteTime           string `json:"quoteTime"`
	QuoteDate           string `json:"quoteDate"`
}

type ReactionType struct {
	ReactionMessageId string `json:"reactionMessageId"`
	ReactionContext   string `json:"reactionContext"`
	ReactionClientId  string `json:"reactionClientId"`
}

type TypingPayload struct {
	PayloadType int    `json:"payloadType"`
	ClientDbId  string `json:"clientDbId"`
	IsTyping    bool   `json:"isTyping"`
}

func Connection(app *tview.Application) error {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	c, _, err := websocket.DefaultDialer.Dial("ws://192.168.178.69:5588/chat", nil)
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
				AddNewEncryptedMessageToChatView(&messagePayload.MessageType.MessageContext)

			case 5:
				var payloadB TypingPayload
				if err := json.Unmarshal(message, &payloadB); err != nil {
					fmt.Println("Error parsing typingPayload:", err)
				}
				AddNewPlainMessageToChatView(&payloadB.ClientDbId)

			default:
				fmt.Println("Unbekannter payloadType:", msg.PayloadType)
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