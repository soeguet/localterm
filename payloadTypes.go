package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

type App struct {
	ui   *tview.Application
	conn *websocket.Conn
}

type MessageListPayload struct {
	PayloadType int              `json:"payloadType"`
	MessageList []MessagePayload `json:"messageList"`
}

type GenericMessage struct {
	PayloadType int             `json:"payloadType"`
	Payload     json.RawMessage `json:"payload"`
}

type MessagePayload struct {
	PayloadType  int             `json:"payloadType"`
	MessageType  MessageType     `json:"messageType"`
	ClientType   ClientType      `json:"clientType"`
	QuoteType    *QuoteType      `json:"quoteType"`
	ReactionType *[]ReactionType `json:"reactionType"`
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

type MessageListRequestPayload struct {
	PayloadType int `json:"payloadType"`
}

type ClientListRequestPayload struct {
	PayloadType int `json:"payloadType"`
}

type QuoteType struct {
	QuoteDbId           string `json:"quoteDbId"`
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

type ReactionPayload struct {
	PayloadType       int    `json:"payloadType"`
	ReactionDbId      string `json:"reactionDbId"`
	ReactionMessageId string `json:"reactionMessageId"`
	ReactionContext   string `json:"reactionContext"`
	ReactionClientId  string `json:"reactionClientId"`
}

type TypingPayload struct {
	PayloadType int    `json:"payloadType"`
	ClientDbId  string `json:"clientDbId"`
	IsTyping    bool   `json:"isTyping"`
}

type AuthenticationPayload struct {
	PayloadType    int    `json:"payloadType"`
	ClientUsername string `json:"clientUsername"`
	ClientDbId     string `json:"clientDbId"`
}

type ClientUpdatePayload struct {
	PayloadType        int    `json:"payloadType"`
	ClientDbId         string `json:"clientDbId"`
	ClientUsername     string `json:"clientUsername"`
	ClientColor        string `json:"clientColor"`
	ClientProfileImage string `json:"clientProfileImage"`
}
