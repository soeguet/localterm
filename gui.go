package main

import (
	"encoding/base64"
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
	chatView *tview.TextView
)

func createChatView(app *App) *tview.TextView {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.ui.Draw()
		})

	return textView
}

func AddNewMessageViaMessagePayload(payload *MessagePayload) {

	decodedString := DecodeBase64ToString(payload.MessageType.MessageContext)

	payloadUsername := GetUsernameForId(payload.ClientType.ClientDbId)

	if payloadUsername == "" {
		payloadUsername = "Unknown"
	}

	usernameColor := GetClientColor(payload.ClientType.ClientDbId)

	fmt.Fprintf(chatView, " [-]%s - ["+usernameColor+"]%s:[-] %s\n", payload.MessageType.MessageTime, payloadUsername,
		decodedString)

	chatView.ScrollToEnd()
}

func AddNewEncryptedMessageToChatView(customMessage *string) {

	decodedString, err := base64.StdEncoding.DecodeString(*customMessage)
	if err != nil {
		fmt.Println("error decoding base64 string:", err)
		return
	}
	fmt.Fprintf(chatView, " [red]%s\n", decodedString)
	chatView.ScrollToEnd()
}

func AddNewPlainMessageToChatView(customMessage *string) {

	fmt.Fprintf(chatView, " [red]%s\n", *customMessage)
	chatView.ScrollToEnd()
}

func AddNewPlainByteToChatView(customMessage *[]byte) {

	fmt.Fprintf(chatView, " [red]%s\n", *customMessage)
	chatView.ScrollToEnd()
}

func AddNewMessageToChatView(payload MessagePayload) {
	encodedmessage := DecodeBase64ToString(payload.MessageType.MessageContext)
	time := payload.MessageType.MessageTime
	sender := payload.ClientType.ClientDbId

	fmt.Fprintf(chatView, " %s - [yellow]%s:[-] %s\n", time, sender, encodedmessage)
	chatView.ScrollToEnd()
}

func createInputField(app *App) *tview.InputField {

	inputField := tview.NewInputField().
		SetLabel("Nachricht: ")

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			// addNewMessageToChatView(inputField.GetText())
			var abc = inputField.GetText()
			//AddNewPlainMessageToChatView(&abc)
			sendMessagePayloadToWebsocket(app.conn, &abc)
			inputField.SetText("")
		}
	})

	return inputField
}

func Gui(app *App) error {

	chatView = createChatView(app)
	inputField := createInputField(app)

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(chatView, 0, 1, false).
		AddItem(inputField, 3, 1, true)

	app.ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	if err := app.ui.SetRoot(flex,
		true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}