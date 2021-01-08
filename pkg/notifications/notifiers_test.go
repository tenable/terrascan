package notifications

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/accurics/terrascan/pkg/config"
	"github.com/accurics/terrascan/pkg/notifications/webhook"
)

func TestNewNotifier(t *testing.T) {

	table := []struct {
		name     string
		nType    string
		wantType Notifier
		wantErr  error
	}{
		{
			name:     "valid notifier",
			nType:    "webhook",
			wantType: &webhook.Webhook{},
			wantErr:  nil,
		},
		{
			name:    "invalid notifier",
			nType:   "notthere",
			wantErr: errNotifierNotSupported,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotErr := NewNotifier(tt.nType)
			if !reflect.DeepEqual(gotType, tt.wantType) {
				t.Errorf("got: '%v', want: '%v'", gotType, tt.wantType)
			}
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("incorrect error; got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}

func TestNewNotifiers(t *testing.T) {

	table := []struct {
		name       string
		configFile string
		wantErr    error
	}{
		{
			name:       "config not present",
			configFile: "notthere",
			wantErr:    config.ErrNotPresent,
		},
		{
			name:       "invalid toml",
			configFile: "testdata/invalid.toml",
			wantErr:    fmt.Errorf("(1, 3): was expecting token =, but got \"am\" instead"),
		},
		{
			name:       "key not present",
			configFile: "testdata/nokey.toml",
			wantErr:    ErrTomlKeyNotPresent,
		},
		{
			name:       "invalid notifier",
			configFile: "testdata/invalid-notifier-type.toml",
			wantErr:    errNotifierNotSupported,
		},
		{
			name:       "empty notifier config",
			configFile: "testdata/empty-notifier-config.toml",
			wantErr:    errNotifierTypeNotPresent,
		},
		{
			name:       "invalid notifier config",
			configFile: "testdata/invalid-notifier-config.toml",
			wantErr:    webhook.ErrNilConfigData,
		},
		{
			name:       "valid multiple notifier config",
			configFile: "testdata/valid-notifier-config.toml",
			wantErr:    nil,
		},
		{
			name:    "no config file specified",
			wantErr: nil,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := NewNotifiers(tt.configFile)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("incorrect error; got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
