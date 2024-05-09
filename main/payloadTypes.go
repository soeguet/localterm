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
	PayloadType payloadType      `json:"payloadType"`
	MessageList []messagePayload `json:"messageList"`
}

type genericMessage struct {
	PayloadType payloadType     `json:"payloadType"`
	Payload     json.RawMessage `json:"payload"`
}

type messagePayload struct {
	PayloadType  payloadType     `json:"payloadType"`
	MessageType  messageType     `json:"messageType"`
	ClientType   clientType      `json:"clientType"`
	QuoteType    *quoteType      `json:"quoteType"`
	ReactionType *[]reactionType `json:"reactionType"`
}

type messageType struct {
	MessageDbID    string `json:"messageDbId"`
	MessageContext string `json:"messageContext"`
	MessageTime    string `json:"messageTime"`
	MessageDate    string `json:"messageDate"`
}

type clientType struct {
	ClientDbID string `json:"clientDbId"`
}

type messageListRequestPayload struct {
	PayloadType payloadType `json:"payloadType"`
}

type clientListRequestPayload struct {
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
	PayloadType       payloadType `json:"payloadType"`
	ReactionDbID      string      `json:"reactionDbId"`
	ReactionMessageID string      `json:"reactionMessageId"`
	ReactionContext   string      `json:"reactionContext"`
	ReactionClientID  string      `json:"reactionClientId"`
}

type typingPayload struct {
	PayloadType payloadType `json:"payloadType"`
	ClientDbID  string      `json:"clientDbId"`
	IsTyping    bool        `json:"isTyping"`
}

type authenticationPayload struct {
	PayloadType    payloadType `json:"payloadType"`
	ClientUsername string      `json:"clientUsername"`
	ClientDbID     string      `json:"clientDbId"`
}

type clientUpdatePayload struct {
	PayloadType        payloadType `json:"payloadType"`
	ClientDbID         string      `json:"clientDbId"`
	ClientUsername     string      `json:"clientUsername"`
	ClientColor        string      `json:"clientColor"`
	ClientProfileImage string      `json:"clientProfileImage"`
}
