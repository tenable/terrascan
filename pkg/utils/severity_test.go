/*
	    Copyright (C) 2022 Tenable, Inc.

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
package utils

import (
	"testing"
)

func TestValidateSeverityInput(t *testing.T) {
	type args struct {
		severity string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Invalid severity input",
			args: args{severity: "invalid"},
			want: false,
		},
		{
			name: "Valid severity input",
			args: args{severity: "High"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateSeverityInput(tt.args.severity); got != tt.want {
				t.Errorf("ValidateSeverityInput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckSeverity(t *testing.T) {
	type args struct {
		ruleSeverity    string
		desiredSeverity string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Equal Severity",
			args: args{ruleSeverity: "High", desiredSeverity: "High"},
			want: true,
		},
		{
			name: "desired severity High",
			args: args{ruleSeverity: "Low", desiredSeverity: "High"},
			want: false,
		},
		{
			name: "Different severity: desired severity Low",
			args: args{ruleSeverity: "Low", desiredSeverity: "Low"},
			want: true,
		},
		{
			name: "Different severity: desired severity Medium",
			args: args{ruleSeverity: "Medium", desiredSeverity: "Medium"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckSeverity(tt.args.ruleSeverity, tt.args.desiredSeverity); got != tt.want {
				t.Errorf("CheckSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinSeverityApplicable(t *testing.T) {
	type args struct {
		ruleSeverity string
		minSeverity  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Min Severity Applicable",
			args: args{ruleSeverity: "Low", minSeverity: "Medium"},
			want: true,
		},
		{
			name: "Min Severity Not Applicable",
			args: args{ruleSeverity: "Medium", minSeverity: "Low"},
			want: false,
		},
		{
			name: "Min Severity Applicable with minSeverity medium",
			args: args{ruleSeverity: "High", minSeverity: "Medium"},
			want: false,
		},
		{
			name: "Min Severity Applicable with minSeverity High",
			args: args{ruleSeverity: "Medium", minSeverity: "High"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinSeverityApplicable(tt.args.ruleSeverity, tt.args.minSeverity); got != tt.want {
				t.Errorf("MinSeverityApplicable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMaxSeverityApplicable(t *testing.T) {
	type args struct {
		ruleSeverity string
		maxSeverity  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Max Severity applicable",
			args: args{ruleSeverity: "High", maxSeverity: "Low"},
			want: true,
		},
		{
			name: "Max Severity not applicable",
			args: args{ruleSeverity: "Low", maxSeverity: "High"},
			want: false,
		},
		{
			name: "Max Severity not applicable",
			args: args{ruleSeverity: "Low", maxSeverity: "High"},
			want: false,
		},
		{
			name: "Max Severity applicable with maxSeverity medium",
			args: args{ruleSeverity: "High", maxSeverity: "Medium"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaxSeverityApplicable(tt.args.ruleSeverity, tt.args.maxSeverity); got != tt.want {
				t.Errorf("MaxSeverityApplicable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMinMaxSeverity(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name            string
		args            args
		wantMinSeverity string
		wantMaxSeverity string
	}{
		{
			name:            "no severity",
			args:            args{body: "none"},
			wantMinSeverity: "",
			wantMaxSeverity: "",
		},
		{
			name:            "min severity set to high",
			args:            args{body: "#ts:minseverity=High"},
			wantMinSeverity: "High",
			wantMaxSeverity: "",
		},
		{
			name:            "max severity set to low",
			args:            args{body: "#ts:maxseverity=Low"},
			wantMinSeverity: "",
			wantMaxSeverity: "Low",
		},
		{
			name:            "max severity set to None",
			args:            args{body: "#ts:maxseverity=None"},
			wantMinSeverity: "",
			wantMaxSeverity: "None",
		},
		{
			name:            "max severity set to low and Min severity set to high",
			args:            args{body: "#ts:maxseverity=LOw  #ts:minseverity=hiGh"},
			wantMinSeverity: "hiGh",
			wantMaxSeverity: "LOw",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMinSeverity, gotMaxSeverity := GetMinMaxSeverity(tt.args.body)
			if gotMinSeverity != tt.wantMinSeverity {
				t.Errorf("GetMinMaxSeverity() gotMinSeverity = %v, want %v", gotMinSeverity, tt.wantMinSeverity)
			}
			if gotMaxSeverity != tt.wantMaxSeverity {
				t.Errorf("GetMinMaxSeverity() gotMaxSeverity = %v, want %v", gotMaxSeverity, tt.wantMaxSeverity)
			}
		})
	}
}
