package main

import (
	"errors"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/rivo/tview"
)

type MockNotifier struct {
	Err   error
	Calls []struct {
		Title   string
		Message string
		Icon    string
	}
}

func (m *MockNotifier) Notify(title, message, icon string) error {
	m.Calls = append(m.Calls, struct {
		Title   string
		Message string
		Icon    string
	}{Title: title, Message: message, Icon: icon})
	return m.Err
}

func TestSendNotification(t *testing.T) {
	mockNotifier := &MockNotifier{}
	app := app{notifier: mockNotifier}

	title := "Test Title"
	message := "Test Message"

	err := app.notifier.Notify(title, message, "")
	if err != nil {
		t.Errorf("sendNotification() error = %v, wantErr %v", err, false)
	}
	if len(mockNotifier.Calls) != 1 {
		t.Fatalf("Expected Notify to be called once, got %v calls", len(mockNotifier.Calls))
	}
	if mockNotifier.Calls[0].Title != title || mockNotifier.Calls[0].Message != message {
		t.Errorf("Notify was not called with the correct parameters: got %v", mockNotifier.Calls[0])
	}

	mockNotifier.Err = errors.New("notification failed")
	err = app.notifier.Notify(title, message, "")
	if err == nil {
		t.Errorf("Expected error from sendNotification() but got nil")
	}
	if err.Error() != "notification failed" {
		t.Errorf("Expected error 'notification failed', got %v", err)
	}
}

func TestBeeepNotifier_Notify(t *testing.T) {
	type args struct {
		title   string
		message string
		icon    string
	}
	tests := []struct {
		name    string
		bn      *beeepNotifier
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bn := &beeepNotifier{}
			if err := bn.Notify(tt.args.title, tt.args.message, tt.args.icon); (err != nil) != tt.wantErr {
				t.Errorf("BeeepNotifier.Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_desktopNotification(t *testing.T) {
	type fields struct {
		ui       *tview.Application
		notifier notifier
		conn     *websocket.Conn
	}
	type args struct {
		payload *messagePayload
	}
	tests := []struct {
		fields  fields
		args    args
		name    string
		wantErr bool
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
			if err := app.desktopNotification(tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("App.desktopNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

