// Copyright (c) 2021 the SC authors. All rights reserved. MIT License.

package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/sc-lang/go-sc/scparse"
	"github.com/spf13/cobra"
)

func newValidateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate <file>",
		Short: "Validate an SC file",
		Long: `validate validates an SC file and reports any errors.

If a path to a file is given as an argument, the file is read. Otherwise,
if no arguments are given validate will read from standard input.

Example:

Reading from a file:

    sctool validate input.sc

Reading from standard input:

    cat input.sc | sctool validate
`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) > 1 {
				return fmt.Errorf("expected 1 arg for file name")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var filename string
			var in io.Reader
			if len(args) == 0 {
				// Set filename for error messages
				filename = "<standard input>"
				in = cmd.InOrStdin()
			} else {
				filename = args[0]
			}
			data, err := readData(filename, in)
			if err != nil {
				return err
			}

			_, err = scparse.Parse(data)
			if err == nil {
				return nil
			}
			var parseErr *scparse.Error
			if !errors.As(err, &parseErr) {
				return fmt.Errorf("unexpected error occurred: %v", err)
			}
			pos := parseErr.Pos
			return fmt.Errorf("syntax error in %s at %d:%d: %s", filename, pos.Line, pos.Column, parseErr.Context)
		},
	}
	return cmd
}

func readData(filename string, in io.Reader) ([]byte, error) {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		in = f
	}
	return io.ReadAll(in)
}
