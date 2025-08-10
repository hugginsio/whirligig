// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package command

import (
	"log"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List local themes",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("todo")
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
