// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package render

import (
	"bytes"

	"github.com/yuin/goldmark"
)

type MarkdownEngine struct{}

func NewMarkdownEngine() *MarkdownEngine {
	return &MarkdownEngine{}
}

func (m *MarkdownEngine) Render(content []byte) ([]byte, error) {
	// TODO: front matter
	// TODO: goldmark config

	var out bytes.Buffer
	if err := goldmark.Convert(content, &out); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
