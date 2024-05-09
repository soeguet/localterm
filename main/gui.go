// main package
package main

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"github.com/rivo/tview"
)

const (
	margin   = "            "
	message  = " Message: "
	quote    = "   Quote: "
	reaction = "Reaction: "
	setting  = "Settings: "
)

var (
	textLabel  string
	chatView   *tview.TextView
	flex       tview.Flex
	typingView *tview.TextView
	inputField *tview.InputField
)

func (app *app) getTypingViewTextLabel() string {
	return typingView.GetLabel()
}

func newApp(ui *tview.Application, conn *websocket.Conn) (*app, error) {
	// initial request to websocket after handshake
	// asks for all RegisteredUsers in a clientList
	err := authenticateClientAtSocket(conn)

	return &app{
		ui:       ui,
		notifier: &beeepNotifier{},
		conn:     conn,
	}, err
}

func changeTextLabelText(text string) {
	inputField.SetLabel(text)
}

func createTypingView(app *app) *tview.TextView {
	textView := tview.NewTextView().
		SetChangedFunc(func() {
			app.ui.Draw()
		})
	textView.SetTextColor(tcell.ColorGray)

	return textView
}

func (app *app) setTypingLabelText(text string) {
	typingView.SetText(text)
	app.ui.Draw()
}

func createChatView(app *app) *tview.TextView {
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.ui.Draw()
		})

	return textView
}

func addNewMessageToScrollPanel(index *int, payload *messagePayload) {
	messageIndex := fmt.Sprintf("[gray][%03d][-]", *index)
	decodedString, err := decodeBase64ToString(payload.MessageType.MessageContext)
	if err != nil {
		fmt.Println("Error decoding base64 to string:", err)
	}

	payloadUsername := getUsernameForID(payload.ClientType.ClientDbID)

	usernameColor := fmt.Sprintf("[%s]", getClientColor(payload.ClientType.ClientDbID))

	var quote string
	if payload.QuoteType != nil {
		quote = checkForQuote(*payload.QuoteType)
	}

	var reactions string
	if payload.ReactionType != nil {
		reactions = checkForReactions(*payload.ReactionType)
	}

	if _, err = fmt.Fprintf(chatView, "%s%s [-]%s - %s%s:[-] %s %s\n", quote, messageIndex,
		payload.MessageType.MessageTime,
		usernameColor,
		payloadUsername, decodedString, reactions); err != nil {
		fmt.Println("Error writing to chatView:", err)
	}

	chatView.ScrollToEnd()
}

func checkForQuote(quoteType quoteType) string {
	if quoteType.QuoteClientID == "" {
		return ""
	}

	msg, err := decodeBase64ToString(quoteType.QuoteMessageContext)
	if err != nil {
		fmt.Println("Error decoding base64 to string:", err)
	}

	quoteString := fmt.Sprintf("%s[#997275]┌ [%s - %s: %s]\n", margin, quoteType.QuoteTime,
		getUsernameForID(quoteType.QuoteClientID), msg)

	return quoteString
}

// checkForReactions checks for reactions in a given list of reaction types and returns a formatted string representing the reactions.
// If the input is empty, it returns an empty string. The reactions are formatted as "[#8B8000]└ [<reaction1> <reaction2> ...][-]".
// The 'margin' constant represents the indentation for the formatted string.
// The 'ReactionType' struct defines the structure of a reaction and is not included in this documentation.
// The function uses the 'strings.Builder' type to build the formatted string efficiently.
// Returns the formatted string representing the reactions or an empty string if an error occurs.
func checkForReactions(reactionType []reactionType) string {
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
func createInputField(app *app) *tview.InputField {
	customInputField := tview.NewInputField()
	inputField = customInputField.
		SetLabel(message).
		SetLabelColor(tcell.ColorGreenYellow).
		SetFieldBackgroundColor(tcell.ColorBlueViolet).
		SetChangedFunc(func(text string) {
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
				customInputField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
				customInputField.SetFieldTextColor(tcell.ColorDarkOrange)
				changeTextLabelText(quote)
			case 2:
				// reaction
				customInputField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
				customInputField.SetFieldTextColor(tcell.ColorGreen)
				changeTextLabelText(reaction)
			case 3:
				// settings change
				customInputField.SetFieldBackgroundColor(tcell.ColorDarkSlateGray)
				customInputField.SetFieldTextColor(tcell.ColorYellow)
				changeTextLabelText(setting)

			default:
				customInputField.SetFieldBackgroundColor(tcell.ColorBlueViolet)
				customInputField.SetFieldTextColor(tcell.ColorWhite)
				changeTextLabelText(message)
			}
		}).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {

				textInput := customInputField.GetText()

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

					customInputField.SetText("")
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
					customInputField.SetText("")
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

						customInputField.SetText("")

					}
					customInputField.SetText("")
					return
				}

				sendMessagePayloadToWebsocket(app.conn, &textInput)
				customInputField.SetText("")
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

	return customInputField
}

func sendProfileUpdateToWebsocket(conn *websocket.Conn, message *string) {
	// schema: /sc

	// remove the first 4 characters
	trimmedMessage := (*message)[4:]

	checkIfTrimmedMessageIsAHexColor := checkIfHexColor(trimmedMessage)

	thisClient := getThisClient()

	if checkIfTrimmedMessageIsAHexColor {
		profileUpdatePayload := clientUpdatePayload{
			PayloadType:        3,
			ClientDbID:         thisClient.ClientDbID,
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

	return matches != nil
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
	regexPattern := `^/s[n,c][0-9]{3} `
	re := regexp.MustCompile(regexPattern)

	matches := re.FindStringSubmatch(text)
	if matches != nil {
		return 3
	}

	return 0
}

func sendReactionPayloadToWebsocketV2(conn *websocket.Conn, message *string) {
	// schema: /r005

	// grab characters in brackets
	trimmedMessageIndex := (*message)[2:5]
	reactedMessagePayload := getMessageFromCache(atoi(trimmedMessageIndex))
	// remove the first 5 characters
	trimmedMessage := (*message)[6:]

	reactionPayload := reactionPayload{
		PayloadType:       7,
		ReactionDbID:      uuid.New().String(),
		ReactionMessageID: reactedMessagePayload.MessageType.MessageDbID,
		ReactionContext:   trimmedMessage,
		ReactionClientID:  envVars.ID,
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

	reactionPayload := reactionPayload{
		PayloadType:       7,
		ReactionDbID:      uuid.New().String(),
		ReactionMessageID: reactedMessagePayload.MessageType.MessageDbID,
		ReactionContext:   trimmedMessage,
		ReactionClientID:  envVars.ID,
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

	messagePayload := messagePayload{
		PayloadType: 1,
		MessageType: messageType{
			MessageDbID:    "TOBEREMOVED",
			MessageContext: base64.StdEncoding.EncodeToString([]byte(trimmedMessage)),
			MessageTime:    time.Now().Format("15:04"),
			MessageDate:    time.Now().Format("2006-01-02"),
		},
		ClientType: clientType{
			ClientDbID: envVars.ID,
		},
		QuoteType: &quoteType{
			QuoteDbID:           quotedMessagePayload.MessageType.MessageDbID,
			QuoteClientID:       quotedMessagePayload.ClientType.ClientDbID,
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

	messagePayload := messagePayload{
		PayloadType: 1,
		MessageType: messageType{
			MessageDbID:    "TOBEREMOVED",
			MessageContext: base64.StdEncoding.EncodeToString([]byte(trimmedMessage)),
			MessageTime:    time.Now().Format("15:04"),
			MessageDate:    time.Now().Format("2006-01-02"),
		},
		ClientType: clientType{
			ClientDbID: envVars.ID,
		},
		QuoteType: &quoteType{
			QuoteDbID:           quotedMessagePayload.MessageType.MessageDbID,
			QuoteClientID:       quotedMessagePayload.ClientType.ClientDbID,
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

func (app *app) clearChatView() {
	chatView.Clear()
}

func createFlex(app *app) tview.Flex {
	flex := tview.NewFlex()
	inputField := createInputField(app)
	typingView = createTypingView(app)
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(chatView, 0, 1, false)
	flex.AddItem(inputField, 1, 1, true)
	flex.AddItem(typingView, 1, 1, false)

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

func gui(app *app) error {
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
