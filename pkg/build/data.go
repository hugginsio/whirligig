// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package build

import (
	"bufio"
	"bytes"
	"fmt"
	"maps"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hugginsio/whirligig/whirligig"
	"sigs.k8s.io/yaml"
)

// TODO: consider moving to separate `pkg/data` with support for different data files such as JSON, TOML, CSV, etc

func (b *Builder) extractData(site *whirligig.Site) error {
	for _, resource := range site.Resources {
		if err := b.loadCompanionData(&resource.File); err != nil {
			return fmt.Errorf("failed to extract data for resource %s: %w", resource.File.Name, err)
		}

		if path.Ext(resource.Name) == ".md" {
			if err := b.loadFrontMatter(resource); err != nil {
				return fmt.Errorf("failed to extract data for resource %s: %w", resource.Name, err)
			}
		}

		if err := b.commonPropertyOverrides(&resource.File); err != nil {
			return err
		}

		// TODO: Resource Data overrides
	}

	for _, file := range site.Files {
		if err := b.loadCompanionData(file); err != nil {
			return fmt.Errorf("failed to extract data for file %s: %w", file.Name, err)
		}

		if err := b.commonPropertyOverrides(file); err != nil {
			return err
		}
	}

	return nil
}

// Load the companion data for a File, e.g. `_styles.yaml`.
func (b *Builder) loadCompanionData(file *whirligig.File) error {
	yamlPath := filepath.Join(b.whirligig.SourcePath, file.Path, "_"+file.Basename+".yaml")

	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		// Try .yml extension as well
		yamlPath = filepath.Join(b.whirligig.SourcePath, file.Path, "_"+file.Basename+".yml")
		if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
			return nil
		}
	}

	yamlContent, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file %s: %w", yamlPath, err)
	}

	var data map[string]any
	if err := yaml.Unmarshal(yamlContent, &data); err != nil {
		return fmt.Errorf("failed to parse YAML file %s: %w", yamlPath, err)
	}

	if file.Data == nil {
		file.Data = make(map[string]any)
	}

	maps.Copy(file.Data, data)
	return nil
}

// Load the frontmatter from a Markdown document.
func (b *Builder) loadFrontMatter(resource *whirligig.Resource) error {
	content, err := resource.Content(b.whirligig.SourcePath)
	if err != nil {
		return fmt.Errorf("failed to read resource content: %w", err)
	}

	frontmatter, err := ParseFrontmatter(content)
	if err != nil {
		return fmt.Errorf("failed to parse front matter: %w", err)
	}

	if resource.Data == nil {
		resource.Data = make(map[string]any)
	}

	if frontmatter == nil {
		return nil
	}

	maps.Copy(resource.Data, frontmatter)

	// TODO: move to separate method later

	if title, ok := frontmatter["title"].(string); ok {
		resource.Title = title
		delete(resource.Data, "title")
	}

	if excerpt, ok := frontmatter["excerpt"].(string); ok {
		resource.Excerpt = excerpt
		delete(resource.Data, "excerpt")
	}

	if url, ok := frontmatter["url"].(string); ok {
		resource.Url = url
		delete(resource.Data, "url")
	}

	return nil
}

// Extract frontmatter from the content of a Markdown document.
func ParseFrontmatter(content []byte) (map[string]any, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))

	// Bail out if content does not start with front matter delimiter
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, nil
	}

	var frontMatterLines []string
	var contentStart int
	lineCount := 1

	for scanner.Scan() {
		line := scanner.Text()
		lineCount++

		if strings.TrimSpace(line) == "---" {
			contentStart = lineCount
			break
		}

		frontMatterLines = append(frontMatterLines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}

	// If we didn't find closing ---, treat as no front matter
	if contentStart == 0 {
		return nil, nil
	}

	frontMatterYAML := strings.Join(frontMatterLines, "\n")
	var frontMatter map[string]any

	if len(frontMatterYAML) > 0 {
		if err := yaml.Unmarshal([]byte(frontMatterYAML), &frontMatter); err != nil {
			return nil, fmt.Errorf("invalid front matter YAML: %w", err)
		}
	} else {
		frontMatter = make(map[string]any)
	}

	return frontMatter, nil
}

func (b *Builder) commonPropertyOverrides(file *whirligig.File) error {
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
