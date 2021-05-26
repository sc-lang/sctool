// Copyright (c) 2021 the SC authors. All rights reserved. MIT License.

package main

import (
	"os"

	"github.com/sc-lang/sctool/cmd"
)

func main() {
	c := cmd.New()
	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
