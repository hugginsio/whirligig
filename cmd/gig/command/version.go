// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package command

import (
	"fmt"

	"github.com/hugginsio/whirligig/internal/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version.GetVersionInfo())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
