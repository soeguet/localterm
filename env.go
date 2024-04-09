package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

var (
	clientList ClientList
	envVars    = EnvVars{
		Username: os.Getenv("LOCALCHAT_USERNAME"),
		IP:       os.Getenv("LOCALCHAT_IP"),
		Port:     os.Getenv("LOCALCHAT_PORT"),
		Os:       runtime.GOOS,
		Id:       setClientId(),
	}
)

type Client struct {
	ClientDbId         string `json:"clientDbId"`
	ClientUsername     string `json:"clientUsername"`
	ClientColor        string `json:"clientColor"`
	ClientProfileImage string `json:"clientProfileImage"`
}

type ClientList struct {
	Clients []Client `json:"clients"`
}

type EnvVars struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Os       string `json:"os"`
	Id       string `json:"id"`
}

func setClientId() string {

	// if dev=true environment variable is set, use a random id
	if os.Getenv("DEV") == "true" {
		return uuid.New().String()
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error retrieving home path: %v", err)
	}

	idFilePath := filepath.Join(homeDir, ".localchat", "id", "id.txt")

	if err := os.MkdirAll(filepath.Dir(idFilePath), 0700); err != nil {
		log.Fatalf("error creating folder: %v", err)
	}

	if _, err := os.Stat(idFilePath); os.IsNotExist(err) {
		// id file missing -> generate new id
		newID := uuid.New().String()

		// save id in file
		if err := os.WriteFile(idFilePath, []byte(newID), 0600); err != nil {
			log.Fatalf("error saving the id: %v", err)
		}

		log.Printf("new id generated and saved: %s", newID)
		return newID
	} else {
		// id exists -> read id from file
		id, err := os.ReadFile(idFilePath)
		if err != nil {
			log.Fatalf("error reading id: %v", err)
		}

		log.Printf("id was read from file: %s", string(id))
		return string(id)
	}
}

func GetThisClientId() string {
	return envVars.Id
}
