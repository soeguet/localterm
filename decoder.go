package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
)

func DecodeJson(jsonString []byte) interface{} {

	var basePayload struct {
		PayloadType int `json:"payloadType"`
	}

	err := json.Unmarshal(jsonString, &basePayload)
	if err != nil {
		log.Fatal(err)
	}

	switch basePayload.PayloadType {
	case 1:
		var payload MessagePayload
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
		var payload MessageListPayload
		err := json.Unmarshal(jsonString, &payload)
		if err != nil {
			log.Fatal(err)
		}
		return payload

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
