// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

// Package renderer provides interfaces for rendering Resources.
package render

type Engine interface {
	Render([]byte) ([]byte, error)
}
