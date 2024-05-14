// main package
package main

import (
	"fmt"
	"log"

	"github.com/gen2brain/beeep"
	"github.com/rivo/tview"
)

// notifier is an interface that defines the behavior of a notification service.
type notifier interface {
	// Notify sends a notification with the specified title, message, and icon.
	// It returns an error if the notification fails to be sent.
	Notify(title, message, icon string) error
}

// beeepNotifier is a type that represents a notifier using the beeep package.
type beeepNotifier struct{}

// Notify sends a notification with the specified title and message using the beeep package. Icons are ignored right now.
func (bn *beeepNotifier) Notify(title, message, _ string) error {
	return beeep.Notify(title, message, "")
}

// desktopNotification sends a desktop notification with the given payload.
// It decodes the base64 encoded message context, retrieves the username from the cache,
// and then calls the notifier to send the notification.
func (app *app) desktopNotification(payload *messagePayload) error {

	decodedString, err := decodeBase64ToString(payload.MessageType.MessageContext)
	if err != nil {
		fmt.Println("Error decoding base64 to string:", err)
	}

	userFromCache := getUsernameFromCache(payload.ClientType.ClientDbID)

	return app.notifier.Notify("message from "+userFromCache, decodedString, "")
}

// createApp creates and initializes a new instance of the app struct.
// It sets up the user interface, establishes a connection, and returns the app.
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
