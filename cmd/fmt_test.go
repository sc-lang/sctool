// Copyright (c) 2021 the SC authors. All rights reserved. MIT License.

package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sc-lang/sctool/cmd"
)

func TestFmt(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		in       io.Reader
		wantFile string
		wantStr  string
	}{
		{
			name:     "formatted file",
			args:     []string{filepath.FromSlash("testdata/valid/basic.sc")},
			wantFile: filepath.FromSlash("testdata/valid/basic.sc.golden"),
		},
		{
			name:     "formatted file with comments",
			args:     []string{filepath.FromSlash("testdata/valid/comments.sc")},
			wantFile: filepath.FromSlash("testdata/valid/comments.sc.golden"),
		},
		{
			name: "format from stdin",
			in: strings.NewReader(`{
				foo: [1
				2,
				null,],},`),
			wantStr: `{
  foo: [
    1
    2
    null
  ]
}
`,
		},
		{
			name: "list unformatted",
			args: []string{"--list", filepath.FromSlash("testdata/valid")},
			wantStr: fmt.Sprintf(
				"%s\n%s\n",
				filepath.FromSlash("testdata/valid/basic.sc"),
				filepath.FromSlash("testdata/valid/comments.sc"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cmd.New()
			if tt.in != nil {
				cmd.SetIn(tt.in)
			}
			args := append([]string{"fmt"}, tt.args...)
			cmd.SetArgs(args)
			var stdout bytes.Buffer
			cmd.SetOut(&stdout)

			err := cmd.Execute()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := stdout.String()
			if tt.wantStr != "" {
				if got != tt.wantStr {
					t.Errorf("got formatted output\n%s\nwant\n%s", got, tt.wantStr)
				}
				return
			}

			data, err := os.ReadFile(tt.wantFile)
			if err != nil {
				t.Fatalf("failed to read %s: %v", tt.wantFile, err)
			}
			want := string(data)
			if got != want {
				t.Errorf("got formatted file\n%s\nwant\n%s", got, want)
			}
		})
	}
}
