package notifications

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tenable/terrascan/pkg/config"
	"github.com/tenable/terrascan/pkg/notifications/webhook"
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

	testDataDir := "testdata"

	table := []struct {
		name       string
		configFile string
		wantErr    error
	}{
		{
			name:       "config not present",
			configFile: "notthere",
			wantErr:    ErrNotificationNotPresent,
		},
		{
			name:       "invalid toml",
			configFile: filepath.Join(testDataDir, "invalid.toml"),
			wantErr:    ErrNotificationNotPresent,
		},
		{
			name:       "key not present",
			configFile: filepath.Join(testDataDir, "nokey.toml"),
			wantErr:    ErrNotificationNotPresent,
		},
		{
			name:       "invalid notifier",
			configFile: filepath.Join(testDataDir, "invalid-notifier-type.toml"),
			wantErr:    errNotifierNotSupported,
		},
		{
			name:       "empty notifier config",
			configFile: filepath.Join(testDataDir, "empty-notifier-config.toml"),
			wantErr:    errNotifierTypeNotPresent,
		},
		{
			name:       "invalid notifier config",
			configFile: filepath.Join(testDataDir, "invalid-notifier-config.toml"),
			wantErr:    webhook.ErrNilConfigData,
		},
		{
			name:       "valid multiple notifier config",
			configFile: filepath.Join(testDataDir, "valid-notifier-config.toml"),
			wantErr:    nil,
		},
		{
			name:    "no config file specified",
			wantErr: ErrNotificationNotPresent,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			config.LoadGlobalConfig(tt.configFile)
			_, gotErr := NewNotifiers()
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("incorrect error; got: '%v', want: '%v'", gotErr, tt.wantErr)
			}
		})
	}
}
