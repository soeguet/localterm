package main

import (
	"encoding/base64"
	"fmt"
	"strings"

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

	usernameColor := fmt.Sprintf("[%s]", GetClientColor(payload.ClientType.ClientDbId))

	var quote string
	if payload.QuoteType != nil {

		quote = checkForQuote(*payload.QuoteType)
	}

	var reactions string
	if payload.ReactionType != nil {
		reactions = checkForReactions(*payload.ReactionType)
	}

	fmt.Fprintf(chatView, "%s [-]%s - %s%s:[-] %s %s\n", quote, payload.MessageType.MessageTime, usernameColor, payloadUsername, decodedString, reactions)

	chatView.ScrollToEnd()
}

func checkForQuote(quoteType QuoteType) string {

	if quoteType.QuoteClientId == "" {
		return ""
	}

	quoteString := fmt.Sprintf("       [gray]┌ [%s - %s: %s]\n", quoteType.QuoteTime, GetUsernameForId(quoteType.QuoteClientId), DecodeBase64ToString(quoteType.QuoteMessageContext))

	return quoteString
}

func checkForReactions(reactionType []ReactionType) string {
	if len(reactionType) == 0 {
		return ""
	}

	var reactions strings.Builder
	reactions.WriteString("\n       [yellow]└ [")
	for _, reaction := range reactionType {
		fmt.Fprintf(&reactions, " %s", reaction.ReactionContext)
	}
	reactions.WriteString("][-]")

	return reactions.String()
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

	fmt.Fprintf(chatView, " [yellow]%s\n", "asd")
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

	inputField := tview.NewInputField()
	inputField.SetLabel("Message: ")
	inputField.SetLabelColor(tcell.ColorGreenYellow)
	inputField.SetFieldBackgroundColor(tcell.ColorBlueViolet)

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			var textInputField = inputField.GetText()
			sendMessagePayloadToWebsocket(app.conn, &textInputField)
			inputField.SetText("")
		}
	})

	return inputField
}

func (app *App) ClearChatView() {
	chatView.Clear()
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
