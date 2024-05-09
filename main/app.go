// main package
package main

import (
	"fmt"
	"log"

	"github.com/gen2brain/beeep"
	"github.com/rivo/tview"
)

type notifier interface {
	Notify(title, message, icon string) error
}

type beeepNotifier struct{}

func (bn *beeepNotifier) Notify(title, message, _ string) error {
	return beeep.Notify(title, message, "")
}

func (app *app) desktopNotification(payload *messagePayload) error {
	decodedString, err := decodeBase64ToString(payload.MessageType.MessageContext)
	if err != nil {
		fmt.Println("Error decoding base64 to string:", err)
	}

	return app.notifier.Notify("message from "+payload.ClientType.ClientDbID, decodedString, "")
}

func createApp() *app {
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
