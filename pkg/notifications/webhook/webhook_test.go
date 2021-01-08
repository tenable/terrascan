/*
    Copyright (C) 2020 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package webhook

import "testing"

func TestWebhookInit(t *testing.T) {
	testURL := "testURL"
	testToken := "testToken"

	type args struct {
		config interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		assert  bool
		url     string
		token   string
	}{
		{
			name: "nil config",
			args: args{
				config: nil,
			},
			wantErr: true,
		},
		{
			name: "valid webhook config data",
			args: args{
				config: map[string]interface{}{
					"url":   testURL,
					"token": testToken,
				},
			},
			assert: true,
			url:    testURL,
			token:  testToken,
		},
		{
			name: "invalid webhook config data",
			args: args{
				config: struct {
					url   string
					token string
				}{
					url:   testURL,
					token: testToken,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Webhook{}
			if err := w.Init(tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("Webhook.Init() got error: %v, wantErr: %v", err, tt.wantErr)
			}
			if tt.assert {
				if w.URL != tt.url {
					t.Errorf("Webhook.Init() got url: %v, want url: %v", w.URL, tt.url)
				}

				if w.Token != tt.token {
					t.Errorf("Webhook.Init() got token: %v, want token: %v", w.Token, tt.token)
				}
			}
		})
	}
}
