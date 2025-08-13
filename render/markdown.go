// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package render

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

type MarkdownEngine struct{}

func NewMarkdownEngine() *MarkdownEngine {
	return &MarkdownEngine{}
}

func (m *MarkdownEngine) Render(content []byte) ([]byte, error) {
	// https://github.com/yuin/goldmark?tab=readme-ov-file#extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			&frontmatter.Extender{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(),
	)

	var out bytes.Buffer
	if err := md.Convert(content, &out); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
