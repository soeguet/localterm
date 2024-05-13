// main package
package main

import (
	"reflect"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

// Define TestCheckForReactions function
func TestCheckForReactions(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		output string
		input  []reactionType
	}{
		{
			name:   "Empty Input",
			input:  []reactionType{},
			output: "",
		},
		{
			name:   "Single Reaction",
			input:  []reactionType{{ReactionContext: "smile"}},
			output: "\n" + margin + "[#8B8000]└ [ smile ][-]",
		},
		{
			name:   "Multiple Reactions",
			input:  []reactionType{{ReactionContext: "smile"}, {ReactionContext: "clap"}, {ReactionContext: "laugh"}},
			output: "\n" + margin + "[#8B8000]└ [ smile clap laugh ][-]",
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkForReactions(tt.input); got != tt.output {
				t.Errorf("checkForReactions() = %v, want %v", got, tt.output)
			}
		})
	}
}

func TestEvalTextInChatView(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "NoMatches",
			text:     "[090] abc ",
			expected: 0,
		},
		{
			name:     "SingleGT",
			text:     "[123] > ",
			expected: 1,
		},
		{
			name:     "DoubleGT",
			text:     "[234] >> ",
			expected: 2,
		},
		{
			name:     "DifferentNumber",
			text:     "[345] >> ",
			expected: 2,
		},
		{
			name:     "InvalidString",
			text:     "invalid text",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evalTextInChatView(tt.text); got != tt.expected {
				t.Errorf("evalTextInChatView() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCheckIfHexColor(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid hex color six digits", "#a3b2c1", true},
		{"valid hex color three digits", "#abc", true},
		{"valid hex color mix cases", "#AaBbCc", true},
		{"invalid hex color not a hex digit", "#xyz", false},
		{"invalid hex color empty", "#", false},
		{"invalid hex color no hash", "abcdef", false},
		{"invalid hex color five digits", "#abcd0", false},
		{"invalid hex color four digits", "#abcd", false},
		{"invalid hex color two digits", "#af", false},
		{"invalid hex color one digit", "#a", false},
		{"invalid hex color seven digits", "#abcdefg", false},
		{"invalid hex color eight digits", "#12345678", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if res := checkIfHexColor(tc.input); res != tc.expected {
				t.Fatalf("expected %v, but got %v", tc.expected, res)
			}
		})
	}
}

func TestEvalTextInChatViewV2(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{
			name: "regex match 'q'",
			text: "/q123 hello",
			want: 1,
		},
		{
			name: "regex match 'r'",
			text: "/r999 hello",
			want: 2,
		},
		{
			name: "regex misMatch",
			text: "/p189 what's up",
			want: 0,
		},
		{
			name: "empty input",
			text: "",
			want: 0,
		},
		{
			name: "regex match without following space",
			text: "/r435hello",
			want: 0,
		},
		{
			name: "regex match but no digit",
			text: "/r hello",
			want: 0,
		},
		{
			name: "text without slash",
			text: "hello world",
			want: 0,
		},
		{
			name: "text with too many digits",
			text: "/q1234 hello",
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evalTextInChatViewV2(tt.text); got != tt.want {
				t.Errorf("evalTextInChatViewV2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvalTextInChatViewV3(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{name: "MatchFound", text: "/sn233 Hello", want: 3},
		{name: "MatchNotFound", text: "/shello", want: 0},
		{name: "EmptyString", text: "", want: 0},
		{name: "NumericContent", text: "/s123456", want: 0},
		{name: "SpecialCharacterContent", text: "/s@#!$%^&*(", want: 0},
		{name: "WhitespaceContent", text: " ", want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evalTextInChatViewV3(tt.text); got != tt.want {
				t.Errorf("evalTextInChatViewV3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAtoi(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "PositiveNumber",
			input:    "123",
			expected: 123,
		},
		{
			name:     "NegativeNumber",
			input:    "-123",
			expected: -123,
		},
		{
			name:     "Zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "InvalidInput",
			input:    "abc",
			expected: 0,
		},
		{
			name:     "EmptyInput",
			input:    "",
			expected: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if output := atoi(tc.input); output != tc.expected {
				t.Errorf("atoi(%v) = %v; want %v", tc.input, output, tc.expected)
			}
		})
	}
}

func Test_newApp(t *testing.T) {
	type args struct {
		ui   *tview.Application
		conn *websocket.Conn
	}
	tests := []struct {
		args    args
		want    *app
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newApp(tt.args.ui, tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("newApp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newApp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_changeTextLabelText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changeTextLabelText(tt.args.text)
		})
	}
}

func Test_createTypingView(t *testing.T) {
	type args struct {
		app *app
	}
	tests := []struct {
		args args
		want *tview.TextView
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createTypingView(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTypingView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_setTypingLabelText(t *testing.T) {
}

func Test_createChatView(t *testing.T) {
	type args struct {
		app *app
	}
	tests := []struct {
		args args
		want *tview.TextView
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createChatView(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createChatView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addNewMessageToScrollPanel(t *testing.T) {
	type args struct {
		index   *int
		payload *messagePayload
	}
	tests := []struct {
		args args
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addNewMessageToScrollPanel(tt.args.index, tt.args.payload)
		})
	}
}

func Test_checkForQuote(t *testing.T) {
	type args struct {
		quoteType quoteType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkForQuote(tt.args.quoteType); got != tt.want {
				t.Errorf("checkForQuote() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checkForReactions(t *testing.T) {
	type args struct {
		reactionType []reactionType
	}
	tests := []struct {
		name string
		want string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkForReactions(tt.args.reactionType); got != tt.want {
				t.Errorf("checkForReactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evalTextInChatView(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evalTextInChatView(tt.args.text); got != tt.want {
				t.Errorf("evalTextInChatView() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createInputField(t *testing.T) {
	type args struct {
		app *app
	}
	tests := []struct {
		args args
		want *tview.InputField
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createInputField(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createInputField() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendProfileUpdateToWebsocket(t *testing.T) {
	type args struct {
		conn    *websocket.Conn
		message *string
	}
	tests := []struct {
		args args
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendProfileUpdateToWebsocket(tt.args.conn, tt.args.message)
		})
	}
}

func Test_checkIfHexColor(t *testing.T) {
	type args struct {
		trimmedMessage string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIfHexColor(tt.args.trimmedMessage); got != tt.want {
				t.Errorf("checkIfHexColor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evalTextInChatViewV2(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "invalid, return 0",
			args: args{
				text: "hello",
			},
			want: 0,
		},
		{
			name: "valid q, return 1",
			args: args{
				text: "/q123 hello",
			},
			want: 1,
		},
		{
			name: "valid r, return 2",
			args: args{
				text: "/r999 hello",
			},
			want: 2,
		},
		{
			name: "invalid q, less than 3 digits, return 0",
			args: args{
				text: "/q12 hello",
			},
			want: 0,
		},
		{
			name: "invalid r, less than 3 digits, return 2",
			args: args{
				text: "/r99 hello",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evalTextInChatViewV2(tt.args.text); got != tt.want {
				t.Errorf("evalTextInChatViewV2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_evalTextInChatViewV3(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "MatchFound",
			args: args{
				text: "/sn333 Hello",
			},
			want: 3,
		},
		{
			name: "MatchNotFound",
			args: args{
				text: "/shello",
			},
			want: 0,
		},
		{
			name: "MatchNotFound",
			args: args{
				text: "/s22",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evalTextInChatViewV3(tt.args.text); got != tt.want {
				t.Errorf("evalTextInChatViewV3() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendReactionPayloadToWebsocketV2(t *testing.T) {
	type args struct {
		conn    *websocket.Conn
		message *string
	}
	tests := []struct {
		args args
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendReactionPayloadToWebsocketV2(tt.args.conn, tt.args.message)
		})
	}
}

func Test_sendReactionPayloadToWebsocket(t *testing.T) {
	type args struct {
		conn    *websocket.Conn
		message *string
	}
	tests := []struct {
		args args
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendReactionPayloadToWebsocket(tt.args.conn, tt.args.message)
		})
	}
}

func Test_sendQuotedMessagePayloadToWebsocketV2(t *testing.T) {
	type args struct {
		conn    *websocket.Conn
		message *string
	}
	tests := []struct {
		args args
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendQuotedMessagePayloadToWebsocketV2(tt.args.conn, tt.args.message)
		})
	}
}

func Test_sendQuotedMessagePayloadToWebsocket(t *testing.T) {
	type args struct {
		conn    *websocket.Conn
		message *string
	}
	tests := []struct {
		args args
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendQuotedMessagePayloadToWebsocket(tt.args.conn, tt.args.message)
		})
	}
}

func Test_atoi(t *testing.T) {
	type args struct {
		index string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := atoi(tt.args.index); got != tt.want {
				t.Errorf("atoi() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApp_clearChatView(t *testing.T) {
	type fields struct {
		ui       *tview.Application
		notifier notifier
		conn     *websocket.Conn
	}
	tests := []struct {
		fields fields
		name   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &app{
				ui:       tt.fields.ui,
				notifier: tt.fields.notifier,
				conn:     tt.fields.conn,
			}
			app.clearChatView()
		})
	}
}

func Test_createFlex(t *testing.T) {
	type args struct {
		app *app
	}
	tests := []struct {
		name string
		args args
		want tview.Flex
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createFlex(tt.args.app); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFlex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_gui(t *testing.T) {
	type args struct {
		app *app
	}
	tests := []struct {
		args    args
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := gui(tt.args.app); (err != nil) != tt.wantErr {
				t.Errorf("gui() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
