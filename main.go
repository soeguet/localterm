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

	go func() {
		err := Connection(app)
		if err != nil {
			log.Fatal(err)
		}
	}()

	if err := Gui(app); err != nil {
		log.Fatal(err)
	}
}