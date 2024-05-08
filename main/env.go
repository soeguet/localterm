package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/google/uuid"
)

func init() {
	setClientId()
}

var (
	mutex               sync.Mutex
	clientUsernameCache = make(map[string]string)
	clientColorCache    = make(map[string]string)
	messageCache        = make(map[int]MessagePayload)
	clientList          ClientList
	thisClient          Client
	envVars             = EnvVars{
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

func getEnvUsername() string {
	if envVars.Username == "" {
		return "Unknown"
	}
	return envVars.Username
}

func getEnvIP() string {
	if envVars.IP == "" {
		return "localhost"
	}
	return envVars.IP
}

func getEnvPort() string {
	if envVars.Port == "" {
		return "8080"
	}
	return envVars.Port
}

func getThisClient() Client {
	return thisClient
}

func resetMessageCache() {
	mutex.Lock()
	defer mutex.Unlock()

	messageCache = make(map[int]MessagePayload)
}
func appendMessageToCache(message MessagePayload) (index int) {
	mutex.Lock()
	defer mutex.Unlock()

	cacheSize := len(messageCache)
	index = cacheSize
	messageCache[index] = message

	return index
}

func getMessageFromCache(index int) MessagePayload {
	mutex.Lock()
	defer mutex.Unlock()

	return messageCache[index]
}

func addUsernameToCache(clientId string, username string) {
	clientUsernameCache[clientId] = username
}

func getUsernameFromCache(clientId string) string {
	mutex.Lock()
	defer mutex.Unlock()

	return clientUsernameCache[clientId]
}

func resetUsernameCache() {
	clientUsernameCache = make(map[string]string)
}

func setClientList(newClientList *ClientList) {
	mutex.Lock()
	defer mutex.Unlock()

	clientList.Clients = newClientList.Clients
	resetUsernameCache()
	resetColorCache()
	cacheThisClient()
}

func cacheThisClient() {
	for _, v := range clientList.Clients {
		if v.ClientDbId == envVars.Id {
			thisClient = v
			return
		}
	}
}

func resetColorCache() {
	clientColorCache = make(map[string]string)
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

	if err := os.MkdirAll(filepath.Dir(idFilePath), 0o700); err != nil {
		log.Fatalf("error creating folder: %v", err)
	}

	if _, err := os.Stat(idFilePath); os.IsNotExist(err) {
		// id file missing -> generate new id
		newID := uuid.New().String()

		// save id in file
		if err := os.WriteFile(idFilePath, []byte(newID), 0o600); err != nil {
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

// getClientColor returns the color of the client with the given client id
// return yellow if the client did not choose a color yet
func getClientColor(clientId string) string {
	mutex.Lock()
	defer mutex.Unlock()

	value, exists := clientColorCache[clientId]

	// if the color is already in the cache, return it; all good
	if exists {
		return value
	}

	// if the color is not in the cache, search for it in the client list and add it to the cache
	for _, v := range clientList.Clients {
		if v.ClientDbId == clientId && v.ClientColor != "" {
			addClientColorToCache(clientId, v.ClientColor)
			return v.ClientColor
		}
	}

	// if the color is not in the cache and not in the client list, return nil
	return "yellow"
}

func addClientColorToCache(id string, color string) {
	clientColorCache[id] = color
}

func getUsernameForId(clientId string) string {
	mutex.Lock()
	defer mutex.Unlock()

	value, exists := clientUsernameCache[clientId]

	// if the username is already in the cache, return it; all good
	if exists {
		return value
	}

	// if the username is not in the cache, search for it in the client list and add it to the cache
	for _, v := range clientList.Clients {
		if v.ClientDbId == clientId {
			addUsernameToCache(clientId, v.ClientUsername)
			return v.ClientUsername
		}
	}

	// if the username is not in the cache and not in the client list, return "Unknown"
	return "Unknown"
}

func getThisClientId() *string {
	return &envVars.Id
}
