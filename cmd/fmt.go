// Copyright (c) 2021 the SC authors. All rights reserved. MIT License.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/sc-lang/go-sc/scparse"
	"github.com/spf13/cobra"
)

type fmtOptions struct {
	list          bool
	write         bool
	errorOnChange bool
}

func newFmtCmd() *cobra.Command {
	var opts fmtOptions
	cmd := &cobra.Command{
		Use:   "fmt [paths...]",
		Short: "Format SC files",
		Long: `fmt formats SC files.

If no path is provided, it will read from standard input. If a provided path is
a directory, it will recursively search for .sc files and format them.

By default, fmt writes the formatted result to standard output. The --write
flag can be used to format files in-place (i.e. overwrite them).

The --list flag can be used to list files that require formatting.

fmt will exit with status code 0 if formatting is successful and status code 1
if an error occurs. The --error-on-change flag can be used to have fmt exit with
status code 2 if any files require formatting.

Example:

Format file and write to standard output:

    sctool fmt input.sc

Reading from standard input:

    cat input.sc | sctool fmt

Format file in-place:

    sctool fmt --write input.sc

List files that require formatting:

    sctool fmt --list configs/
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			changed, err := runFmt(cmd, args, opts)
			if err != nil {
				return err
			}
			if opts.errorOnChange && len(changed) > 0 {
				fmt.Fprintln(cmd.ErrOrStderr(), "file(s) were changed:", strings.Join(changed, ", "))
				// Explicitly exit with code 2 to indicate change
				os.Exit(2)
			}
			return nil
		},
	}
	cmd.Flags().BoolVarP(&opts.list, "list", "l", false, "list files whose formatting differs from sctool's")
	cmd.Flags().BoolVarP(&opts.write, "write", "w", false, "write result to source files instead of stdout")
	cmd.Flags().BoolVarP(&opts.errorOnChange, "error-on-change", "e", false, "exit with status code 2 if any files are changed during formatting")
	return cmd
}

func runFmt(cmd *cobra.Command, args []string, opts fmtOptions) ([]string, error) {
	stdout := cmd.OutOrStdout()
	var changedFiles []string

	// Read from stdin if no files provided
	if len(args) == 0 {
		if opts.list {
			return nil, fmt.Errorf("--list not allowed on stdin")
		}
		if opts.write {
			return nil, fmt.Errorf("--write not allowed on stdin")
		}
		changed, err := formatFile("", cmd.InOrStdin(), stdout, opts)
		if changed {
			changedFiles = append(changedFiles, "<standard input>")
		}
		return changedFiles, err
	}

	for _, path := range args {
		info, err := os.Stat(path)
		if err != nil {
			return changedFiles, err
		}
		if !info.IsDir() {
			changed, err := formatFile(path, nil, stdout, opts)
			if changed {
				changedFiles = append(changedFiles, path)
			}
			if err != nil {
				return changedFiles, fmt.Errorf("unabled to format %s: %v", path, err)
			}
			continue
		}
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || !strings.HasSuffix(info.Name(), ".sc") {
				return nil
			}
			changed, err := formatFile(path, nil, stdout, opts)
			if changed {
				changedFiles = append(changedFiles, path)
			}
			if err != nil {
				err = fmt.Errorf("unabled to format %s: %v", path, err)
			}
			return err
		})
		if err != nil {
			return changedFiles, err
		}
	}
	return changedFiles, nil
}

func formatFile(filename string, in io.Reader, out io.Writer, opts fmtOptions) (bool, error) {
	if in == nil {
		f, err := os.Open(filename)
		if err != nil {
			return false, err
		}
		defer f.Close()
		in = f
	}

	src, err := io.ReadAll(in)
	if err != nil {
		return false, err
	}
	n, err := scparse.Parse(src)
	if err != nil {
		return false, err
	}
	n = processDictionary(n)
	res := scparse.Format(n)

	changed := false
	if !bytes.Equal(src, res) {
		changed = true
		if opts.list {
			fmt.Fprintln(out, filename)
		}
		if opts.write {
			var perms os.FileMode
			if info, err := os.Stat(filename); err == nil {
				perms = info.Mode() & os.ModePerm
			}
			err = os.WriteFile(filename, res, perms)
			if err != nil {
				return changed, err
			}
		}
	}
	if !opts.list && !opts.write {
		_, err = out.Write(res)
	}
	return changed, err
}

func processDictionary(n *scparse.DictionaryNode) *scparse.DictionaryNode {
	modified := false
	var newMembers []*scparse.MemberNode
	for _, m := range n.Members {
		if m.Key.Type() == scparse.NodeIdentifier {
			newMember := m
			v := processValue(m.Value)
			if v != m.Value {
				modified = true
				newMember = &scparse.MemberNode{Pos: m.Pos, CommentGroup: m.CommentGroup, Key: m.Key, Value: v}
			}
			newMembers = append(newMembers, newMember)
			continue
		}

		// See if we can unquote key
		keyStr := m.Key.KeyString()
		needsQuote := false
		for i, r := range keyStr {
			if r != '_' && !unicode.IsLetter(r) && (i == 0 || !unicode.IsDigit(r)) {
				needsQuote = true
				break
			}
		}
		v := processValue(m.Value)
		if needsQuote && v == m.Value {
			newMembers = append(newMembers, m)
			continue
		}

		modified = true
		newKey := m.Key
		if !needsQuote {
			newKey = &scparse.IdentifierNode{Pos: m.Key.Position(), CommentGroup: *m.Key.Comments(), Name: keyStr}
		}
		newMember := &scparse.MemberNode{Pos: m.Pos, CommentGroup: m.CommentGroup, Key: newKey, Value: v}
		newMembers = append(newMembers, newMember)
	}
	if modified {
		return &scparse.DictionaryNode{Pos: n.Pos, CommentGroup: n.CommentGroup, Members: newMembers}
	}
	return n
}

func processValue(n scparse.ValueNode) scparse.ValueNode {
	switch n := n.(type) {
	case *scparse.DictionaryNode:
		return processDictionary(n)
	case *scparse.ListNode:
		modified := false
		var newElements []scparse.ValueNode
		for _, e := range n.Elements {
			v := processValue(e)
			if v != e {
				modified = true
			}
			newElements = append(newElements, v)
		}
		if modified {
			return &scparse.ListNode{Pos: n.Pos, CommentGroup: n.CommentGroup, Elements: newElements}
		}
	}
	return n
}
