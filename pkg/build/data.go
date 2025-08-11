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

	"github.com/hugginsio/whirligig/whirligig"
	"sigs.k8s.io/yaml"
)

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

		// TODO: common Data overrides
		// TODO: Resource Data overrides
	}

	for _, file := range site.Files {
		if err := b.loadCompanionData(file); err != nil {
			return fmt.Errorf("failed to extract data for file %s: %w", file.Name, err)
		}

		// TODO: common Data overrides
	}

	return nil
}

func (b *Builder) loadCompanionData(file *whirligig.File) error {
	yamlPath := filepath.Join(b.sourcePath, file.Path, "_"+file.Basename+".yaml")

	if _, err := os.Stat(yamlPath); os.IsNotExist(err) {
		// Try .yml extension as well
		yamlPath = filepath.Join(b.sourcePath, file.Path, file.Basename+".yml")
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

func (b *Builder) loadFrontMatter(resource *whirligig.Resource) error {
	content, err := resource.Content(b.sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read resource content: %w", err)
	}

	frontMatter, _, err := parseFrontMatter(content)
	if err != nil {
		return fmt.Errorf("failed to parse front matter: %w", err)
	}

	if resource.Data == nil {
		resource.Data = make(map[string]any)
	}

	if frontMatter == nil {
		return nil
	}

	maps.Copy(resource.Data, frontMatter)

	// TODO: move to separate method later

	if title, ok := frontMatter["title"].(string); ok {
		resource.Title = title
		delete(resource.Data, "title")
	}

	if excerpt, ok := frontMatter["excerpt"].(string); ok {
		resource.Excerpt = excerpt
		delete(resource.Data, "excerpt")
	}

	if url, ok := frontMatter["url"].(string); ok {
		resource.Url = url
		delete(resource.Data, "url")
	}

	return nil
}

// parseFrontMatter splits Markdown frontmatter from its content.
func parseFrontMatter(content []byte) (map[string]any, []byte, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))

	// Bail out if content does not start with front matter delimiter
	if !scanner.Scan() || strings.TrimSpace(scanner.Text()) != "---" {
		return nil, content, nil
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
		return nil, nil, fmt.Errorf("error reading content: %w", err)
	}

	// If we didn't find closing ---, treat as no front matter
	if contentStart == 0 {
		return nil, content, nil
	}

	frontMatterYAML := strings.Join(frontMatterLines, "\n")
	var frontMatter map[string]any

	if len(frontMatterYAML) > 0 {
		if err := yaml.Unmarshal([]byte(frontMatterYAML), &frontMatter); err != nil {
			return nil, nil, fmt.Errorf("invalid front matter YAML: %w", err)
		}
	} else {
		frontMatter = make(map[string]any)
	}

	// TODO: need to see if this is even necessary with goldmark renderer
	// Extract content after front matter
	lines := strings.Split(string(content), "\n")
	if contentStart < len(lines) {
		remainingContent := strings.Join(lines[contentStart:], "\n")
		return frontMatter, []byte(remainingContent), nil
	}

	return frontMatter, []byte{}, nil
}
