// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package command

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hugginsio/whirligig/pkg/build"
	"github.com/hugginsio/whirligig/whirligig"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Build your site",
	Run: func(cmd *cobra.Command, args []string) {
		builder := build.New(SourcePath, whirligig.Configuration{})
		if err := builder.Prepare(); err != nil {
			log.Fatalln(err)
		}

		// Print site as pretty JSON for dev purposes
		siteJSON, err := json.MarshalIndent(builder.GetSite(), "", "  ")
		if err != nil {
			log.Fatalf("failed to marshal site to JSON: %v", err)
		}

		fmt.Println(string(siteJSON))

		// TODO: dry run flag that bails out after printing site data

		if err := builder.Build(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
