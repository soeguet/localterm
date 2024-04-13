package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gorilla/websocket"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	chatView *tview.TextView
	margin   = "            "
)

func NewApp(ui *tview.Application, localChatIp string, localChatPort string) (*App, error) {
	url := fmt.Sprintf("ws://%s:%s/chat", localChatIp, localChatPort)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	// initial request to websocket after handshake
	// asks for all RegisteredUsers in a clientList
	authenticateClientAtSocket(conn)
	// retrieveLast100Messages(conn)

	return &App{
		ui:   ui,
		conn: conn,
	}, nil
}

func createChatView(app *App) *tview.TextView {

	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.ui.Draw()
		})

	return textView
}

func AddNewMessageViaMessagePayload(index *uint16, payload *MessagePayload) {

	messageIndex := fmt.Sprintf("[gray][%03d][-]", *index)
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

	_, err := fmt.Fprintf(chatView, "%s%s [-]%s - %s%s:[-] %s %s\n", quote, messageIndex,
		payload.MessageType.MessageTime,
		usernameColor,
		payloadUsername, decodedString, reactions)

	if err != nil {
		fmt.Println("Error writing to chatView:", err)
	}

	chatView.ScrollToEnd()
}

func checkForQuote(quoteType QuoteType) string {

	if quoteType.QuoteClientId == "" {
		return ""
	}

	quoteString := fmt.Sprintf("%s[#997275]┌ [%s - %s: %s]\n", margin, quoteType.QuoteTime,
		GetUsernameForId(quoteType.QuoteClientId), DecodeBase64ToString(quoteType.QuoteMessageContext))

	return quoteString
}

func checkForReactions(reactionType []ReactionType) string {
	if len(reactionType) == 0 {
		return ""
	}

	var reactions strings.Builder

	reactions.WriteString("\n" + margin + "[#8B8000]└ [")
	for _, reaction := range reactionType {
		_, err := fmt.Fprintf(&reactions, " %s", reaction.ReactionContext)
		if err != nil {
			return ""
		}
	}
	reactions.WriteString(" ][-]")

	return reactions.String()
}

func AddNewPlainMessageToChatView(customMessage *string) {

	fmt.Fprintf(chatView, " [yellow]%s\n", "asd")
	chatView.ScrollToEnd()
}

func createInputField(app *App) *tview.InputField {

	inputField := tview.NewInputField()
	inputField.SetLabel("Message: ")
	inputField.SetLabelColor(tcell.ColorGreenYellow)
	inputField.SetFieldBackgroundColor(tcell.ColorBlueViolet)

	regexPattern := `^\[([0-9]{3})\] (>{1,2}) `
	re := regexp.MustCompile(regexPattern)

	inputField.SetChangedFunc(func(text string) {
		matches := re.FindStringSubmatch(text)
		if matches != nil {
			switch matches[2] {
			case ">":
				// quote
				inputField.SetFieldTextColor(tcell.ColorDarkOrange)
			case ">>":
				// reaction
				inputField.SetFieldTextColor(tcell.ColorGreen)
			default:
				inputField.SetFieldTextColor(tcell.ColorWhite)
			}
		} else {
			inputField.SetFieldTextColor(tcell.ColorWhite)
		}
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			textInputField := inputField.GetText()
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

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(chatView, 0, 1, false)
	flex.AddItem(inputField, 3, 1, true)

	app.ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	if err := app.ui.SetRoot(flex,
		true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}