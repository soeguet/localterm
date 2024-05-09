// main package
package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/google/uuid"
)

func init() {
	setClientID()
}

var (
	mutex               sync.Mutex
	clientUsernameCache = make(map[string]string)
	clientColorCache    = make(map[string]string)
	messageCache        = make(map[int]messagePayload)
	typingClientCache   = []string{}
	clientList          clientListStruct
	thisClient          client
	envVars             = envVarsStruct{
		Username: os.Getenv("LOCALCHAT_USERNAME"),
		IP:       os.Getenv("LOCALCHAT_IP"),
		Port:     os.Getenv("LOCALCHAT_PORT"),
		Os:       runtime.GOOS,
		ID:       setClientID(),
	}
)

type client struct {
	ClientDbID         string `json:"clientDbId"`
	ClientUsername     string `json:"clientUsername"`
	ClientColor        string `json:"clientColor"`
	ClientProfileImage string `json:"clientProfileImage"`
}

type clientListStruct struct {
	Clients []client `json:"clients"`
}

type envVarsStruct struct {
	Username string `json:"username"`
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Os       string `json:"os"`
	ID       string `json:"id"`
}

func addTypingClient(clientID string) {
	mutex.Lock()
	defer mutex.Unlock()
	typingClientCache = append(typingClientCache, clientID)
}

func removeTypingClient(clientID string) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, v := range typingClientCache {
		if v == clientID {
			typingClientCache = append(typingClientCache[:i], typingClientCache[i+1:]...)
			return
		}
	}
}

func generateTypingString() string {
	// this freezes the UI
	// mutex.Lock()
	// defer mutex.Unlock()

	if len(typingClientCache) == 0 {
		return ""
	}
	if len(typingClientCache) == 1 {
		return getUsernameForID(typingClientCache[0]) + " is typing..."
	}

	// if there are multiple clients typing
	var builder strings.Builder
	length := len(typingClientCache)

	for i, v := range typingClientCache {
		username := getUsernameForID(v)
		if i == length-1 && i != 0 {
			builder.WriteString("and ")
		}
		builder.WriteString(username)
		if i < length-1 {
			if i == length-2 {
				builder.WriteString(" ")
			} else {
				builder.WriteString(", ")
			}
		}
	}

	builder.WriteString(" are typing...")
	return builder.String()
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

func getThisClient() client {
	return thisClient
}

func resetMessageCache() {
	mutex.Lock()
	defer mutex.Unlock()

	messageCache = make(map[int]messagePayload)
}

func appendMessageToCache(message messagePayload) (index int) {
	mutex.Lock()
	defer mutex.Unlock()

	cacheSize := len(messageCache)
	index = cacheSize
	messageCache[index] = message

	return index
}

func getMessageFromCache(index int) messagePayload {
	mutex.Lock()
	defer mutex.Unlock()

	return messageCache[index]
}

func addUsernameToCache(clientID string, username string) {
	clientUsernameCache[clientID] = username
}

func getUsernameFromCache(clientID string) string {
	mutex.Lock()
	defer mutex.Unlock()

	return clientUsernameCache[clientID]
}

func resetUsernameCache() {
	clientUsernameCache = make(map[string]string)
}

func setClientList(newClientList *clientListStruct) {
	mutex.Lock()
	defer mutex.Unlock()

	clientList.Clients = newClientList.Clients
	resetUsernameCache()
	resetColorCache()
	cacheThisClient()
}

func cacheThisClient() {
	for _, v := range clientList.Clients {
		if v.ClientDbID == envVars.ID {
			thisClient = v
			return
		}
	}
}

func resetColorCache() {
	clientColorCache = make(map[string]string)
}

func setClientID() string {
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
	}

	// id exists -> read id from file
	id, err := os.ReadFile(idFilePath)
	if err != nil {
		log.Fatalf("error reading id: %v", err)
	}

	log.Printf("id was read from file: %s", string(id))
	return string(id)
}

// getClientColor returns the color of the client with the given client id
// return yellow if the client did not choose a color yet
func getClientColor(clientID string) string {
	mutex.Lock()
	defer mutex.Unlock()

	value, exists := clientColorCache[clientID]

	// if the color is already in the cache, return it; all good
	if exists {
		return value
	}

	// if the color is not in the cache, search for it in the client list and add it to the cache
	for _, v := range clientList.Clients {
		if v.ClientDbID == clientID && v.ClientColor != "" {
			addClientColorToCache(clientID, v.ClientColor)
			return v.ClientColor
		}
	}

	// if the color is not in the cache and not in the client list, return nil
	return "yellow"
}

func addClientColorToCache(id string, color string) {
	clientColorCache[id] = color
}

func getUsernameForID(clientID string) string {
	mutex.Lock()
	defer mutex.Unlock()

	value, exists := clientUsernameCache[clientID]

	// if the username is already in the cache, return it; all good
	if exists {
		return value
	}

	// if the username is not in the cache, search for it in the client list and add it to the cache
	for _, v := range clientList.Clients {
		if v.ClientDbID == clientID {
			addUsernameToCache(clientID, v.ClientUsername)
			return v.ClientUsername
		}
	}

	// if the username is not in the cache and not in the client list, return "Unknown"
	return "Unknown"
}

func getThisClientID() *string {
	return &envVars.ID
}
