// Copyright (c) 2021 the SC authors. All rights reserved. MIT License.

package cmd_test

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sc-lang/sctool/cmd"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		path string
		in   io.Reader
		want string
	}{
		{
			name: "valid file",
			path: "testdata/valid/comments.sc",
		},
		{
			name: "invalid file",
			path: "testdata/invalid/missing_value.sc",
			want: fmt.Sprintf(
				"syntax error in %s at 3:1: unexpected <}> in value",
				filepath.FromSlash("testdata/invalid/missing_value.sc"),
			),
		},
		{
			name: "valid stdin",
			in:   strings.NewReader(`{ foo: true, bar: null }`),
		},
		{
			name: "invalid stdin",
			in:   strings.NewReader(`{ foo: [ }`),
			want: "syntax error in <standard input> at 1:10: unexpected <}> in value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd.New()
			if tt.in != nil {
				cmd.SetIn(tt.in)
			}
			args := []string{"validate"}
			if tt.path != "" {
				args = append(args, filepath.FromSlash(tt.path))
			}
			cmd.SetArgs(args)

			err := cmd.Execute()
			if err == nil {
				if tt.want != "" {
					t.Error("want non-nil error")
				}
				return
			}
			if err.Error() != tt.want {
				t.Errorf("got error\n\t%s\nwant\n\t%s", err, tt.want)
			}
		})
	}
}
