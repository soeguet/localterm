package main

import (
	"testing"
)

// Define TestCheckForReactions function
func TestCheckForReactions(t *testing.T) {
	// Define test cases
	tests := []struct {
		name   string
		input  []ReactionType
		output string
	}{
		{
			name:   "Empty Input",
			input:  []ReactionType{},
			output: "",
		},
		{
			name:   "Single Reaction",
			input:  []ReactionType{{ReactionContext: "smile"}},
			output: "\n" + margin + "[#8B8000]└ [ smile ][-]",
		},
		{
			name:   "Multiple Reactions",
			input:  []ReactionType{{ReactionContext: "smile"}, {ReactionContext: "clap"}, {ReactionContext: "laugh"}},
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
		{name: "MatchFound", text: "/sn Hello", want: 3},
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
