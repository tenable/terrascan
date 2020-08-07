package runtime

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/notifications"
)

// MockNotifier mocks notifications.Notifier interface
type MockNotifier struct {
	err error
}

var (
	errMockNotifier = fmt.Errorf("mock notification error")
)

func (m MockNotifier) Init(config interface{}) error {
	return m.err
}

func (m MockNotifier) SendNotification(config interface{}) error {
	return m.err
}

func TestSendNotifications(t *testing.T) {

	table := []struct {
		name     string
		executor Executor
		wantErr  error
	}{
		{
			name: "no notifier error",
			executor: Executor{
				notifiers: []notifications.Notifier{&MockNotifier{err: nil}},
			},
			wantErr: nil,
		},
		{
			name: "no notifier error",
			executor: Executor{
				notifiers: []notifications.Notifier{&MockNotifier{err: errMockNotifier}},
			},
			wantErr: errMockNotifier,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.executor.SendNotifications("some data")
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("incorrect error; got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
