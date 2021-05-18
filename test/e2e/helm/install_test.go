package helm

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestHelmInstall(t *testing.T) {

	table := []struct {
		description string
		helmOptions HelmOptions
		wantErr     error
	}{
		{
			description: "default dry-run test",
			helmOptions: HelmOptions{
				HelmBin:   "helm",
				Name:      "terrascan",
				DryRun:    true,
				Debug:     true,
				ChartPath: os.Getenv("DEFAULT_CHART_PATH"),
				Namespace: "test",
			},
			wantErr: nil,
		},

		{
			description: "secret key file empty error",
			helmOptions: HelmOptions{
				HelmBin:   "helm",
				Name:      "terrascan",
				DryRun:    true,
				Debug:     true,
				ChartPath: os.Getenv("DEFAULT_CHART_PATH"),
				SetFlags: map[string]string{
					"secrets.tlsKeyFilePath": "secret.key",
				},
				Namespace: "test",
			},
			wantErr: fmt.Errorf(`unknown object type "nil" in Secret.data.tls.crt`),
		},
	}

	for _, tt := range table {
		t.Run(tt.description, func(t *testing.T) {
			_, err := tt.helmOptions.Install(context.Background())
			if tt.wantErr != nil && err != nil && strings.Contains(err.Error(), tt.wantErr.Error()) {
				return
			}

			if tt.wantErr != err {
				t.Errorf("unexpected error; got: '%v', want: '%v'", err, tt.wantErr)
			}

		})
	}
}
