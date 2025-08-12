// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package data

import (
	"fmt"
	"maps"
	"path/filepath"

	"github.com/hugginsio/whirligig/whirligig"
)

// EnrichFiles extracts Data for all Files (and Resources) using frontmatter or companion Data files.
func EnrichFiles(site *whirligig.Site, sourcePath string) error {
	for _, resource := range site.Resources {
		if err := enrichFile(sourcePath, &resource.File); err != nil {
			return fmt.Errorf("failed to enrich resource %s: %w", resource.Name, err)
		}

		// TODO: if markdown, load front matter
		// TODO: common property overrides
	}

	return nil
}

func enrichFile(sourcePath string, file *whirligig.File) error {
	filepath := filepath.Join(sourcePath, file.Path)
	data, err := LoadYAML(filepath)
	if err != nil {
		return err
	}

	if file.Data == nil {
		file.Data = make(map[string]any)
	}

	maps.Copy(file.Data, data)

	return nil
}

// EnrichSite extracts Data for the Site as a whole using the data directory.
func EnrichSite(site *whirligig.Site, sourcePath string) error {
	return fmt.Errorf("not implemented")
}
