// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package command

import (
	"github.com/spf13/cobra"
)

var Verbose bool
var rootCmd = &cobra.Command{
	Use: "gig-theme",
}

func Execute() error {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "increased log level")

	return rootCmd.Execute()
}
