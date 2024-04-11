package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"

	"github.com/rivo/tview"
)

func NewApp(ui *tview.Application, localChatIp string, localChatPort string) (*App, error) {
	url := fmt.Sprintf("ws://%s:%s/chat", localChatIp, localChatPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	return &App{
		ui:   ui,
		conn: conn,
	}, nil
}

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