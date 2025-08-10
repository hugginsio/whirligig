// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package main

import (
	"os"

	"github.com/hugginsio/whirligig/cmd/gig/command"
)

func main() {
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
