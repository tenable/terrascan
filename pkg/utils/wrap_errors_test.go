package utils

import (
	"fmt"
	"reflect"
	"testing"
)

func TestWrapError(t *testing.T) {

	mockErr := fmt.Errorf("mock error")

	table := []struct {
		name    string
		err     error
		allErr  error
		wantErr error
	}{
		{
			name:    "empty allErrs",
			allErr:  nil,
			err:     mockErr,
			wantErr: mockErr,
		},
		{
			name:    "empty err",
			err:     nil,
			allErr:  mockErr,
			wantErr: mockErr,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := WrapError(tt.err, tt.allErr)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("incorrect error; got: '%+v', want: '%+v'", gotErr, tt.wantErr)
			}
		})
	}

	t.Run("wrapped error", func(t *testing.T) {
		var (
			err     = fmt.Errorf("mock err")
			allErrs = fmt.Errorf("mock allErrs")
			wantErr = fmt.Errorf("%s: %s", allErrs.Error(), err.Error())
		)
		gotErr := WrapError(err, allErrs)
		if gotErr.Error() != wantErr.Error() {
			t.Errorf("incorrect error: got: '%v', want: '%v'", gotErr, wantErr)
		}
	})
}
