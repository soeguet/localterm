// main package
package main

import (
	"log"
)

func main() {
	app := createApp()

	// Start the connection in a goroutine
	go func() {
		if err := connection(app); err != nil {
			log.Fatal(err)
		}
	}()

	// Start the GUI in the main thread
	if err := gui(app); err != nil {
		log.Fatal(err)
	}
}
