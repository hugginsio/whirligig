// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

// Package build provides tools for building a static site.
package build

import (
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/hugginsio/whirligig/internal/version"
	"github.com/hugginsio/whirligig/pkg/data"
	"github.com/hugginsio/whirligig/pkg/render"
	"github.com/hugginsio/whirligig/whirligig"
)

type Builder struct {
	configuration   *whirligig.Configuration
	destinationPath string
	site            *whirligig.Site
	whirligig       *whirligig.Whirligig
}

// New creates the Builder from the provided configuration.
func New(sourcePath string, configuration whirligig.Configuration) *Builder {
	version := version.GetVersionInfo()
	builder := &Builder{
		configuration: &configuration,
		whirligig: &whirligig.Whirligig{
			SourcePath:  sourcePath,
			Version:     version.GitVersion,
			VersionInfo: &version,
		},
	}

	builder.destinationPath = path.Join(sourcePath, "_site") // TODO: pull from configuration

	return builder
}

// GetSite returns the Site metadata. Returns nil if Prepare has not been called.
func (b *Builder) GetSite() *whirligig.Site {
	return b.site
}

// GetWhirligig returns the Whirligig struct. Returns nil if Prepare has not been called.
func (b *Builder) GetWhirligig() *whirligig.Whirligig {
	return b.whirligig
}

// Clean removes the destination directory and all its children.
func (b *Builder) Clean() error {
	return os.RemoveAll(b.destinationPath)
}

// Prepare walks the specified sourcePath to collect internal Site metadata.
func (b *Builder) Prepare() error {
	b.site = &whirligig.Site{
		Time: time.Now(),
		// TODO: estimate file counts and preallocate slices accordingly
		Resources: make([]*whirligig.Resource, 0),
		Files:     make([]*whirligig.File, 0),
		// TODO: site version
	}

	if err := b.walkSourceDirectory(b.site); err != nil {
		return fmt.Errorf("failed to walk source directory: %w", err)
	}

	if err := data.EnrichFiles(b.site, b.whirligig.SourcePath); err != nil {
		return fmt.Errorf("failed to extract data: %w", err)
	}

	return nil
}

// Build utilizes the Site metadata collected by Prepare to execute the build process and export
// the generated Site to the filesystem.
func (b *Builder) Build() error {
	if b.site == nil {
		return fmt.Errorf("site metadata not prepared. builder.Prepare() is required.")
	}

	if err := b.Clean(); err != nil {
		return fmt.Errorf("failed to clean destination directory: %w", err)
	}

	for _, resource := range b.site.Resources {
		var engine render.Engine

		switch path.Ext(resource.Name) {
		case ".md":
			engine = render.NewMarkdownEngine()
		default:
			// TODO: default engine
			log.Fatalln("default engine not yet implemented")
			engine = nil
		}

		bytes, err := resource.Content(b.whirligig.SourcePath)
		if err != nil {
			return err
		}

		content, err := engine.Render(bytes)
		if err != nil {
			return err
		}

		if err := resource.Write(b.whirligig.SourcePath, b.destinationPath, content); err != nil {
			return err
		}
	}

	// TODO: create gitignore in destinationPath so it is excluded automatically?
	// TODO: doing this concurrently could provide a performance benefit
	for _, file := range b.site.Files {
		file.Copy(b.whirligig.SourcePath, b.destinationPath)
	}

	return nil
}
