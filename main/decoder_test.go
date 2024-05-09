// main package
package main

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBase64ToString(t *testing.T) {
	tests := []struct {
		name          string
		encodedString string
		want          string
		wantErr       bool
	}{
		{
			name:          "empty string",
			encodedString: "",
			want:          "",
			wantErr:       true,
		},
		{
			name:          "valid base64",
			encodedString: base64.StdEncoding.EncodeToString([]byte("hello")),
			want:          "hello",
			wantErr:       false,
		},
		{
			name:          "invalid base64",
			encodedString: "invalid_base64@@",
			want:          "",
			wantErr:       true,
		},
		{
			name:          "decode whitespace",
			encodedString: base64.StdEncoding.EncodeToString([]byte(" ")),
			want:          " ",
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeBase64ToString(tt.encodedString)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
