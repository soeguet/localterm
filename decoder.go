package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

type ClientType struct {
	ClientDbId string `json:"clientDbId"`
}

type MessageType struct {
	MessageDbId    string `json:"messageDbId"`
	MessageContext string `json:"messageContext"`
	MessageTime    string `json:"messageTime"`
	MessageDate    string `json:"messageDate"`
}

type Payload struct {
	PayloadType int         `json:"payloadType"`
	ClientType  ClientType  `json:"clientType"`
	MessageType MessageType `json:"messageType"`
}

func DecodeJson(jsonString []byte) interface{} {

	var basePayload struct {
		PayloadType int `json:"payloadType"`
	}

	err := json.Unmarshal(jsonString, &basePayload)
	if err != nil {
		log.Fatal(err)
	}

	AddNewPlainMessageToChatView("PayloadType: " + string(jsonString))

	switch basePayload.PayloadType {
	case 1:
		var payload Payload
		err := json.Unmarshal(jsonString, &payload)
		if err != nil {
			log.Fatal(err)
		}
		return payload

	case 2:
		var payload ClientList
		err := json.Unmarshal(jsonString, &payload)
		if err != nil {
			log.Fatal(err)
		}
		return payload

	case 4:

	default:
		return -1
	}
}

func DecodeBase64ToString(encodedString string) string {
	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		log.Fatal(err)
	}
	return string(decoded)
}
