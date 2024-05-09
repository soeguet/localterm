// main package
package main

import (
	"encoding/base64"
	"errors"
)

func decodeBase64ToString(encodedString string) (string, error) {
	if encodedString == "" {
		return "", errors.New("input string is empty")
	}

	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
