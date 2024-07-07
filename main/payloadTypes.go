// main package
package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

type app struct {
	ui       *tview.Application
	notifier notifier
	conn     *websocket.Conn
}

type messageListPayload struct {
	MessageList []messagePayload `json:"messageList"`
	PayloadType payloadType      `json:"payloadType"`
}

type genericMessage struct {
	Payload     json.RawMessage `json:"payload"`
	PayloadType payloadType     `json:"payloadType"`
}

type messagePayload struct {
	QuoteType    *quoteType      `json:"quoteType"`
	ReactionType *[]reactionType `json:"reactionType"`
	ImageType    *imageType      `json:"imageType"`
	ClientType   clientType      `json:"clientType"`
	MessageType  messageType     `json:"messageType"`
	PayloadType  payloadType     `json:"payloadType"`
}

type imageType struct {
	ImageDbID string `json:"imageDbId"`
	Type      string `json:"type"`
	Data      string `json:"data"`
}

type messageType struct {
	MessageDbID    string `json:"messageDbId"`
	MessageContext string `json:"messageContext"`
	MessageTime    string `json:"messageTime"`
	MessageDate    string `json:"messageDate"`
	Deleted        bool   `json:"deleted"`
	Edited         bool   `json:"edited"`
}

type clientType struct {
	ClientDbID string `json:"clientDbId"`
}

type messageListRequestPayload struct {
	PayloadType payloadType `json:"payloadType"`
}

type quoteType struct {
	QuoteDbID           string `json:"quoteDbId"`
	QuoteClientID       string `json:"quoteClientId"`
	QuoteMessageContext string `json:"quoteMessageContext"`
	QuoteTime           string `json:"quoteTime"`
	QuoteDate           string `json:"quoteDate"`
}

type reactionType struct {
	ReactionMessageID string `json:"reactionMessageId"`
	ReactionContext   string `json:"reactionContext"`
	ReactionClientID  string `json:"reactionClientId"`
}

type reactionPayload struct {
	ReactionDbID      string      `json:"reactionDbId"`
	ReactionMessageID string      `json:"reactionMessageId"`
	ReactionContext   string      `json:"reactionContext"`
	ReactionClientID  string      `json:"reactionClientId"`
	PayloadType       payloadType `json:"payloadType"`
}

type typingPayload struct {
	ClientDbID  string      `json:"clientDbId"`
	PayloadType payloadType `json:"payloadType"`
	IsTyping    bool        `json:"isTyping"`
}

type authenticationPayload struct {
	ClientUsername string      `json:"clientUsername"`
	ClientDbID     string      `json:"clientDbId"`
	PayloadType    payloadType `json:"payloadType"`
}

type clientUpdatePayload struct {
	ClientDbID         string      `json:"clientDbId"`
	ClientUsername     string      `json:"clientUsername"`
	ClientColor        string      `json:"clientColor"`
	ClientProfileImage string      `json:"clientProfileImage"`
	PayloadType        payloadType `json:"payloadType"`
}
