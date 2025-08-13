// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package data

import (
	"fmt"
	"maps"
	"path"
	"path/filepath"
	"time"

	"github.com/hugginsio/whirligig/whirligig"
)

// EnrichFiles extracts Data for all Files (and Resources) using frontmatter or companion Data files.
func EnrichFiles(site *whirligig.Site, sourcePath string) error {
	for _, resource := range site.Resources {
		if err := enrichFile(sourcePath, &resource.File); err != nil {
			return fmt.Errorf("failed to enrich resource %s: %w", resource.Name, err)
		}

		if path.Ext(resource.Name) == ".md" {
			frontmatter, err := LoadMarkdown(filepath.Join(sourcePath, resource.Path, resource.Name))

			if err != nil {
				return err
			}

			if frontmatter != nil {
				maps.Copy(resource.Data, frontmatter)
			}

			if title, ok := frontmatter["title"].(string); ok {
				resource.Title = title
				delete(resource.Data, "title")
			}
		}

		if err := commonPropertyOverrides(&resource.File); err != nil {
			return err
		}

		if err := resourcePropertyOverrides(resource); err != nil {
			return err
		}
	}

	for _, file := range site.Files {
		if err := enrichFile(sourcePath, file); err != nil {
			return fmt.Errorf("failed to enrich file %s: %w", file.Name, err)
		}

		if err := commonPropertyOverrides(file); err != nil {
			return err
		}
	}

	return nil
}

func enrichFile(sourcePath string, file *whirligig.File) error {
	if file.Data == nil {
		file.Data = make(map[string]any)
	}

	filePathNoExt := filepath.Join(sourcePath, file.Path, "_"+file.Basename)
	data, err := LoadYAML(filePathNoExt + ".yaml")
	if err != nil {
		return err
	}

	maps.Copy(file.Data, data)

	return nil
}

// EnrichSite extracts Data for the Site as a whole using the data directory.
func EnrichSite(site *whirligig.Site, sourcePath string) error {
	return fmt.Errorf("not implemented")
}

func commonPropertyOverrides(file *whirligig.File) error {
	if created, ok := file.Data["created"].(string); ok {
		layouts := []string{
			time.RFC3339,           // "2006-01-02T15:04:05Z07:00"
			"2006-01-02T15:04:05Z", // UTC variant
			"2006-01-02 15:04:05",  // "2023-12-01 15:30:00"
			"2006-01-02",           // "2023-12-01"
			"01/02/2006",           // "12/01/2023"
		}

		var parsedTime time.Time
		var err error

		for _, layout := range layouts {
			if parsedTime, err = time.Parse(layout, created); err == nil {
				file.Created = parsedTime
				break
			}
		}

		if err != nil {
			return fmt.Errorf("Warning: could not parse created timestamp '%s' with any known format: %v\n", created, err)
		}

		delete(file.Data, "created")
	}

	return nil
}

func resourcePropertyOverrides(resource *whirligig.Resource) error {
	if title, ok := resource.Data["title"].(string); ok {
		resource.Title = title

		delete(resource.Data, "title")
	}

	if excerpt, ok := resource.Data["excerpt"].(string); ok {
		resource.Excerpt = excerpt

		delete(resource.Data, "excerpt")
	}

	if url, ok := resource.Data["url"].(string); ok {
		resource.Url = url

		delete(resource.Data, "url")
	}

	return nil
}
