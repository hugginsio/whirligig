// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package command

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

var SourcePath string
var Verbose bool
var rootCmd = &cobra.Command{
	Use:   "gig",
	Short: "Whirligig is a configuration-optional static site generator",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		if SourcePath == "" {
			SourcePath = wd
		}

		if !path.IsAbs(SourcePath) {
			SourcePath = path.Join(wd, path.Clean(SourcePath))
		}

		if _, err := os.Stat(SourcePath); os.IsNotExist(err) {
			return err
		}

		// TODO: set up logger

		return nil
	},
}

func Execute() error {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "increased log level")
	rootCmd.PersistentFlags().StringVarP(&SourcePath, "source", "s", "", "target source directory") // TODO: consider moving this and prerun to build

	return rootCmd.Execute()
}
