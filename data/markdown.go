// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package data

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"sigs.k8s.io/yaml"
)

// LoadMarkdown loads frontmatter from a Markdown document into a hashmap.
func LoadMarkdown(path string) (map[string]any, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

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

	frontmatterYAML := strings.Join(frontMatterLines, "\n")
	var frontmatter map[string]any

	if len(frontmatterYAML) > 0 {
		if err := yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter); err != nil {
			return nil, fmt.Errorf("invalid front matter YAML: %w", err)
		}
	} else {
		frontmatter = make(map[string]any)
	}

	return frontmatter, nil
}
