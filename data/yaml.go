// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"
)

// LoadYAML loads a YAML file into a hashmap. It supports `.yaml` and `.yml` extensions.
func LoadYAML(path string) (map[string]any, error) {
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".yaml" && ext != ".yml" {
		return nil, fmt.Errorf("invalid file extension: %s", ext)
	}

	parent := filepath.Dir(path)
	name := filepath.Base(path[:len(path)-len(ext)])

	if _, err := os.Stat(path); os.IsNotExist(err) {
		path = filepath.Join(parent, name+".yml")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil, nil
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %w", path, err)
	}

	var result map[string]any
	if err := yaml.Unmarshal(content, &result); err != nil {
		return nil, fmt.Errorf("failed to parse YAML file %s: %w", path, err)
	}

	return result, nil
}
