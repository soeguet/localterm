package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func createTextArea() *tview.TextArea {
	textArea := tview.NewTextArea().
		SetPlaceholder("Enter text here...")
	textArea.SetTitle("Text Area").SetBorder(true)
	return textArea
}

var (
	chatView = createChatView()
	appCopy  *tview.Application
)

func createChatView() *tview.TextView {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			appCopy.Draw()
		})

	return textView
}

func AddNewPlainMessageToChatView(customMessage string) {

	fmt.Fprintf(chatView, " [red]%s\n", customMessage)
	chatView.ScrollToEnd()
}

func AddNewMessageToChatView(payload Payload) {
	encodedmessage := DecodeBase64ToString(payload.MessageType.MessageContext)
	time := payload.MessageType.MessageTime
	sender := payload.ClientType.ClientDbId

	fmt.Fprintf(chatView, " %s - [yellow]%s:[-] %s\n", time, sender, encodedmessage)
	chatView.ScrollToEnd()
}

func createInputField() *tview.InputField {

	// Erstelle das InputField, aber setze SetDoneFunc später
	inputField := tview.NewInputField().
		SetLabel("Nachricht: ")

	// Setze nun die Done-Funktion, nachdem inputField bereits definiert ist
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			// addNewMessageToChatView(inputField.GetText())
			AddNewPlainMessageToChatView(inputField.GetText())
			inputField.SetText("")
		}
	})

	return inputField
}

func Gui(app *tview.Application) error {

	appCopy = app
	inputField := createInputField()

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(inputField, 3, 1, true)

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	if err := app.SetRoot(flex,
		true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}
