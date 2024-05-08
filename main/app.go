package main

import (
	"fmt"
	"log"

	"github.com/gen2brain/beeep"
	"github.com/rivo/tview"
)

type Notifier interface {
	Notify(title, message, icon string) error
}

type BeeepNotifier struct{}

func (bn *BeeepNotifier) Notify(title, message, icon string) error {
	return beeep.Notify(title, message, "")
}

func (app *App) desktopNotification(payload *MessagePayload) error {

	decodedString, err := decodeBase64ToString(payload.MessageType.MessageContext)

	if err != nil {
		fmt.Println("Error decoding base64 to string:", err)
	}

	return app.notifier.Notify("message from "+payload.ClientType.ClientDbId, decodedString, "")
}

func createApp() *App {
	ui := tview.NewApplication()

	conn, err := createConnection(getEnvIP(), getEnvPort())
	if err != nil {
		log.Fatalf("Failed to create connection: %v", err)
	}

	app, err := newApp(ui, conn)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	return app
}
