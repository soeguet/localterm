package main

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
type ClientType struct {
	ClientDbId string `json:"clientDbId"`
}

type MessageType struct {
	MessageDbId    string `json:"messageDbId"`
	MessageContext string `json:"messageContext"`
	MessageTime    string `json:"messageTime"`
	MessageDate    string `json:"messageDate"`
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

type MessagePayload struct {
	PayloadType  int           `json:"payloadType"`
	ClientType   ClientType    `json:"clientType"`
	MessageType  MessageType   `json:"messageType"`
	QuoteType    *QuoteType    `json:"quoteType"`
	ReactionType *ReactionType `json:"reactionType"`
}

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
	PayloadType int           `json:"payloadType"`
	MessageList []interface{} `json:"messageList"` // Dies wird spezifisch angepasst
}
