package main

import (
	"errors"
	"testing"
)

type MockNotifier struct {
	Calls []struct {
		Title   string
		Message string
		Icon    string
	}
	Err error
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
	app := App{notifier: mockNotifier}

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
