package main

import (
	"log"

	"github.com/rivo/tview"
)

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
