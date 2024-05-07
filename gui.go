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

const (
	margin = "            "
)

var (
	chatView *tview.TextView
	flex     tview.Flex
)

func newApp(ui *tview.Application, conn *websocket.Conn) (*App, error) {

	// initial request to websocket after handshake
	// asks for all RegisteredUsers in a clientList
	err := authenticateClientAtSocket(conn)

	return &App{
		ui:   ui,
		conn: conn,
	}, err
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

func addNewMessageViaMessagePayload(index *int, payload *MessagePayload) {
	messageIndex := fmt.Sprintf("[gray][%03d][-]", *index)
	decodedString, err := decodeBase64ToString(payload.MessageType.MessageContext)
	if err != nil {
		fmt.Println("Error decoding base64 to string:", err)
	}

	payloadUsername := getUsernameForId(payload.ClientType.ClientDbId)

	usernameColor := fmt.Sprintf("[%s]", getClientColor(payload.ClientType.ClientDbId))

	var quote string
	if payload.QuoteType != nil {
		quote = checkForQuote(*payload.QuoteType)
	}

	var reactions string
	if payload.ReactionType != nil {
		reactions = checkForReactions(*payload.ReactionType)
	}

	_, err = fmt.Fprintf(chatView, "%s%s [-]%s - %s%s:[-] %s %s\n", quote, messageIndex,
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

	msg, err := decodeBase64ToString(quoteType.QuoteMessageContext)
	if err != nil {
		fmt.Println("Error decoding base64 to string:", err)
	}

	quoteString := fmt.Sprintf("%s[#997275]┌ [%s - %s: %s]\n", margin, quoteType.QuoteTime,
		getUsernameForId(quoteType.QuoteClientId), msg)

	return quoteString
}

// checkForReactions checks for reactions in a given list of reaction types and returns a formatted string representing the reactions.
// If the input is empty, it returns an empty string. The reactions are formatted as "[#8B8000]└ [<reaction1> <reaction2> ...][-]".
// The 'margin' constant represents the indentation for the formatted string.
// The 'ReactionType' struct defines the structure of a reaction and is not included in this documentation.
// The function uses the 'strings.Builder' type to build the formatted string efficiently.
// Returns the formatted string representing the reactions or an empty string if an error occurs.
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

// checks for [000] > or [000] >> in the message
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

// createInputField creates a new tview.InputField with a specific label, label color, and field background color.
// It also sets up event handlers for text changes and when the enter key is pressed.
// The input field's field text color is set based on the evaluation of the input text using evalTextInChatView.
// It returns the created input field.
func createInputField(app *App) *tview.InputField {
	inputField := tview.NewInputField()
	inputField.SetLabel("Message: ")
	inputField.SetLabelColor(tcell.ColorGreenYellow)
	inputField.SetFieldBackgroundColor(tcell.ColorBlueViolet)

	inputField.SetChangedFunc(func(text string) {
		textCase := evalTextInChatView(text)

		if textCase == 0 {
			textCase = evalTextInChatViewV2(text)
		}

		if textCase == 0 {
			textCase = evalTextInChatViewV3(text)
		}

		switch textCase {
		case 1:
			// quote
			inputField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
			inputField.SetFieldTextColor(tcell.ColorDarkOrange)
		case 2:
			// reaction
			inputField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
			inputField.SetFieldTextColor(tcell.ColorGreen)
		case 3:
			// settings change
			inputField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
			inputField.SetFieldTextColor(tcell.ColorYellow)
		default:
			inputField.SetFieldBackgroundColor(tcell.ColorBlueViolet)
			inputField.SetFieldTextColor(tcell.ColorWhite)
		}
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {

			textInput := inputField.GetText()

			// check for [000] > or [000] >> in the message
			textCase := evalTextInChatView(textInput)

			if textCase != 0 {
				switch textCase {
				case 1:
					// quote
					sendQuotedMessagePayloadToWebsocket(app.conn, &textInput)
				case 2:
					// reaction
					sendReactionPayloadToWebsocket(app.conn, &textInput)
				default:
					// plain message
					sendMessagePayloadToWebsocket(app.conn, &textInput)
				}

				inputField.SetText("")
				return
			}

			// check for /q or /r followed by three digits and a space
			textCaseV2 := evalTextInChatViewV2(textInput)
			if textCaseV2 != 0 {
				switch textCaseV2 {
				case 1:
					sendQuotedMessagePayloadToWebsocketV2(app.conn, &textInput)
				// quote
				case 2:
					sendReactionPayloadToWebsocketV2(app.conn, &textInput)
				// reaction
				default:
					// plain message
					sendMessagePayloadToWebsocket(app.conn, &textInput)
				}
				inputField.SetText("")
				return
			}

			textCaseV3 := evalTextInChatViewV3(textInput)

			if textCaseV3 != 0 {
				switch textCaseV3 {
				case 4:
					// settings change
					sendProfileUpdateToWebsocket(app.conn, &textInput)
				case 5:
				// settings change
				default:
					// plain message
					sendMessagePayloadToWebsocket(app.conn, &textInput)

					inputField.SetText("")

				}
				inputField.SetText("")
				return
			}

			sendMessagePayloadToWebsocket(app.conn, &textInput)
			inputField.SetText("")
		}
	})

	// inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
	// 	if event.Key() == tcell.KeyF1 {
	// 		fmt.Println("F1 pressed")
	// 		app.ui.SetFocus(modal)
	//
	// 		return nil
	// 	}
	// 	return event
	// })

	return inputField
}

func sendProfileUpdateToWebsocket(conn *websocket.Conn, message *string) {
	// schema: /sc

	// remove the first 4 characters
	trimmedMessage := (*message)[4:]

	checkIfTrimmedMessageIsAHexColor := checkIfHexColor(trimmedMessage)

	thisClient := getThisClient()

	if checkIfTrimmedMessageIsAHexColor {
		profileUpdatePayload := ClientUpdatePayload{
			PayloadType:        3,
			ClientDbId:         thisClient.ClientDbId,
			ClientColor:        trimmedMessage,
			ClientUsername:     thisClient.ClientUsername,
			ClientProfileImage: thisClient.ClientProfileImage,
		}

		err := conn.WriteJSON(profileUpdatePayload)
		if err != nil {
			fmt.Println("Error writing messagePayload:", err)
		}
	}
}

func checkIfHexColor(trimmedMessage string) bool {
	regexPattern := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	re := regexp.MustCompile(regexPattern)

	matches := re.FindStringSubmatch(trimmedMessage)
	if matches != nil {
		return true
	} else {
		return false
	}
}

// checks for /q or /r followed by three digits and a space
func evalTextInChatViewV2(text string) int {
	regexPattern := `^/(q|r)[0-9]{3} `
	re := regexp.MustCompile(regexPattern)

	matches := re.FindStringSubmatch(text)
	if matches != nil {
		switch matches[1] {
		case "q":
			return 1
		case "r":
			return 2
		default:
			return 0
		}
	} else {
		return 0
	}
}

// checks for /s followed by three digits and a space
func evalTextInChatViewV3(text string) int {
	regexPattern := `^/s[n,c] `
	re := regexp.MustCompile(regexPattern)

	matches := re.FindStringSubmatch(text)
	if matches != nil {
		return 3
	} else {
		return 0
	}
}

func sendReactionPayloadToWebsocketV2(conn *websocket.Conn, message *string) {
	// schema: /r005

	// grab characters in brackets
	trimmedMessageIndex := (*message)[2:5]
	reactedMessagePayload := getMessageFromCache(atoi(trimmedMessageIndex))
	// remove the first 5 characters
	trimmedMessage := (*message)[6:]

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

func sendReactionPayloadToWebsocket(conn *websocket.Conn, message *string) {
	// schema: [000] >>

	// grab characters in brackets
	trimmedMessageIndex := (*message)[1:4]
	reactedMessagePayload := getMessageFromCache(atoi(trimmedMessageIndex))
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

func sendQuotedMessagePayloadToWebsocketV2(conn *websocket.Conn, message *string) {
	// schema: /q000

	// grab characters in brackets
	trimmedMessageIndex := (*message)[2:5]
	quotedMessagePayload := getMessageFromCache(atoi(trimmedMessageIndex))
	// remove the first 5 characters
	trimmedMessage := (*message)[5:]

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

func sendQuotedMessagePayloadToWebsocket(conn *websocket.Conn, message *string) {
	// schema: [000] >

	// grab characters in brackets
	trimmedMessageIndex := (*message)[1:4]
	quotedMessagePayload := getMessageFromCache(atoi(trimmedMessageIndex))
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

func (app *App) clearChatView() {
	chatView.Clear()
}

func createFlex(app *App) tview.Flex {
	flex := tview.NewFlex()
	inputField := createInputField(app)
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(chatView, 0, 1, false)
	flex.AddItem(inputField, 3, 1, true)

	return *flex
}

// func createModal(app *App) tview.Modal {
// 	modal := tview.NewModal().
// 		SetText("Do you want to quit the application?").
// 		AddButtons([]string{"Quit", "Cancel"}).
// 		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
// 		}).
// 		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
//
// 			if event.Key() == tcell.KeyEsc {
// 				err := app.ui.SetFocus(&flex)
// 				if err != nil {
// 					return nil
// 				}
// 			}
// 			return event
// 		})
//
// 	return modal
// }

func gui(app *App) error {
	chatView = createChatView(app)
	flex = createFlex(app)
	// modal = createModal(app)

	app.ui.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		return event
	})

	if err := app.ui.SetRoot(&flex,
		true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}

	return nil
}
