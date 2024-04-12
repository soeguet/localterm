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
//type ClientType struct {
//	ClientDbId string `json:"clientDbId"`
//}
//
//type MessageType struct {
//	MessageDbId    string `json:"messageDbId"`
//	MessageContext string `json:"messageContext"`
//	MessageTime    string `json:"messageTime"`
//	MessageDate    string `json:"messageDate"`
//}
//
//type QuoteType struct {
//	QuoteMessageId      string `json:"quoteMessageId"`
//	QuoteClientId       string `json:"quoteClientId"`
//	QuoteMessageContext string `json:"quoteMessageContext"`
//	QuoteTime           string `json:"quoteTime"`
//	QuoteDate           string `json:"quoteDate"`
//}
//
//type ReactionType struct {
//	ReactionMessageId string `json:"reactionMessageId"`
//	ReactionContext   string `json:"reactionContext"`
//	ReactionClientId  string `json:"reactionClientId"`
//}
//
//type MessagePayload struct {
//	PayloadType  int           `json:"payloadType"`
//	ClientType   ClientType    `json:"clientType"`
//	MessageType  MessageType   `json:"messageType"`
//	QuoteType    *QuoteType    `json:"quoteType"`
//	ReactionType *ReactionType `json:"reactionType"`
//}

/**
 * [[ RESULTING TYPE ]]
 * export type MessageListPayload = {
 *    payloadType: PayloadSubType.messageList;
 *    messageList: [
 * *      messageType: {
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
 *    ]
 *
 */
type MessageListPayload struct {
	PayloadType int              `json:"payloadType"`
	MessageList []MessagePayload `json:"messageList"`
}

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

type AuthenticationPayload struct {
	PayloadType    int    `json:"payloadType"`
	ClientUsername string `json:"clientUsername"`
	ClientDbId     string `json:"clientDbId"`
}
