// Copyright (c) 2021 the SC authors. All rights reserved. MIT License.

package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"
)

var version string

func New() *cobra.Command {
	if version == "" {
		if info, available := debug.ReadBuildInfo(); available {
			version = info.Main.Version
		} else {
			version = "master"
		}
	}
	cmd := &cobra.Command{
		Use:          "sctool",
		Version:      version,
		Short:        "sctool is a collection of small utilities for working with SC files.",
		SilenceUsage: true,
	}
	cmd.AddCommand(
		newFmtCmd(),
		newValidateCmd(),
	)
	return cmd
}
