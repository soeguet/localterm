package main

import (
	"log"

	"github.com/rivo/tview"
)

func main() {

	app := tview.NewApplication()

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
