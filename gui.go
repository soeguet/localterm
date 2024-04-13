package main

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
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

func AddNewMessageViaMessagePayload(index *int, payload *MessagePayload) {

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

func evalTextInChatView(text string) int {

	regexPattern := `^\[([0-9]{3})\] (>{1,2}) `
	re := regexp.MustCompile(regexPattern)

	matches := re.FindStringSubmatch(text)
	if matches != nil {
		switch matches[2] {
		case ">":
			// quote
			return 1
		case ">>":
			// reaction
			return 2
		default:
			return 0
		}
	} else {
		return 0
	}
}

func createInputField(app *App) *tview.InputField {

	inputField := tview.NewInputField()
	inputField.SetLabel("Message: ")
	inputField.SetLabelColor(tcell.ColorGreenYellow)
	inputField.SetFieldBackgroundColor(tcell.ColorBlueViolet)

	inputField.SetChangedFunc(func(text string) {
		textCase := evalTextInChatView(text)
		switch textCase {
		case 1:
			// quote
			inputField.SetFieldTextColor(tcell.ColorDarkOrange)
		case 2:
			// reaction
			inputField.SetFieldTextColor(tcell.ColorGreen)
		default:
			inputField.SetFieldTextColor(tcell.ColorWhite)
		}
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			textInputField := inputField.GetText()
			textCase := evalTextInChatView(textInputField)
			switch textCase {
			case 1:
				sendQuotedMessagePayloadToWebsocket(app.conn, &textInputField)
			// quote
			case 2:
				sendReactionPayloadToWebsocket(app.conn, &textInputField)
			// reaction
			default:
				// plain message
				sendMessagePayloadToWebsocket(app.conn, &textInputField)
			}
			inputField.SetText("")
		}
	})

	return inputField
}

func sendReactionPayloadToWebsocket(conn *websocket.Conn, message *string) {

	// schema: [000] >>

	// grab characters in brackets
	trimmedMessageIndex := (*message)[1:4]
	reactedMessagePayload := GetMessageFromCache(atoi(trimmedMessageIndex))
	// remove the first 8 characters
	trimmedMessage := (*message)[8:]

	reactionPayload := ReactionPayload{
		PayloadType:       7,
		ReactionDbId:      uuid.New().String(),
		ReactionMessageId: reactedMessagePayload.MessageType.MessageDbId,
		ReactionContext:   trimmedMessage,
		ReactionClientId:  envVars.Id,
	}

	err := conn.WriteJSON(reactionPayload)
	if err != nil {
		fmt.Println("Error writing messagePayload:", err)
	}

}

func sendQuotedMessagePayloadToWebsocket(conn *websocket.Conn, message *string) {

	// schema: [000] >

	// grab characters in brackets
	trimmedMessageIndex := (*message)[1:4]
	quotedMessagePayload := GetMessageFromCache(atoi(trimmedMessageIndex))
	// remove the first 7 characters
	trimmedMessage := (*message)[7:]

	messagePayload := MessagePayload{
		PayloadType: 1,
		MessageType: MessageType{
			MessageDbId:    "TOBEREMOVED",
			MessageContext: base64.StdEncoding.EncodeToString([]byte(trimmedMessage)),
			MessageTime:    time.Now().Format("15:04"),
			MessageDate:    time.Now().Format("2006-01-02"),
		},
		ClientType: ClientType{
			ClientDbId: envVars.Id,
		},
		QuoteType: &QuoteType{
			QuoteDbId:           quotedMessagePayload.MessageType.MessageDbId,
			QuoteClientId:       quotedMessagePayload.ClientType.ClientDbId,
			QuoteMessageContext: quotedMessagePayload.MessageType.MessageContext,
			QuoteTime:           quotedMessagePayload.MessageType.MessageTime,
			QuoteDate:           quotedMessagePayload.MessageType.MessageDate,
		},
		ReactionType: nil,
	}

	err := conn.WriteJSON(messagePayload)
	if err != nil {
		fmt.Println("Error writing messagePayload:", err)
	}
}

func atoi(index string) int {

	i, err := strconv.Atoi(index)
	if err != nil {
		return 0
	}

	return i
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