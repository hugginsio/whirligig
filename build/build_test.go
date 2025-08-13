// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package build_test

import (
	"testing"

	"github.com/hugginsio/whirligig/build"
	"github.com/hugginsio/whirligig/whirligig"
)

func BenchmarkWalkSourceDirectory(b *testing.B) {
	builder := build.New("../../example", whirligig.Configuration{})

	b.ResetTimer()
	b.ReportAllocs()

	for b.Loop() {
		if err := builder.Prepare(); err != nil {
			b.Fatal(err)
		}

		if err := builder.Build(); err != nil {
			b.Fatal(err)
		}
	}
}
