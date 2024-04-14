package main

import (
	"log"

	"github.com/rivo/tview"
)

func main() {
	ui := tview.NewApplication()
	localChatIp := envVars.IP
	localChatPort := envVars.Port

	app, err := NewApp(ui, localChatIp, localChatPort)
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	// Start the connection in a goroutine
	go func() {
		if err := Connection(app); err != nil {
			log.Fatal(err)
		}
	}()

	// Start the GUI in the main goroutine
	if err := Gui(app); err != nil {
		log.Fatal(err)
	}
}

