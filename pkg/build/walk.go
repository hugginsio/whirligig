// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package build

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/hugginsio/whirligig/whirligig"
)

// Walks the site source directory to build the list of Resources and Files.
func (b *Builder) walkSourceDirectory(site *whirligig.Site) error {
	err := filepath.Walk(b.whirligig.SourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// TODO: should Data files be required to have an underscore?
		if strings.HasPrefix(info.Name(), "_") || strings.HasPrefix(info.Name(), ".") {
			// NOTE: sentinel value, skip the entire directory if bearing exclusion prefix
			if info.IsDir() {
				return filepath.SkipDir
			}

			return nil
		}

		if info.IsDir() {
			return nil
		}

		relativeDir, err := filepath.Rel(b.whirligig.SourcePath, filepath.Dir(path))
		if err != nil {
			return err
		}

		file := whirligig.File{
			Name:     info.Name(),
			Basename: strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
			Extname:  filepath.Ext(info.Name()),
			Created:  b.getCreatedTime(info),
			Modified: info.ModTime(),
			Path:     relativeDir,
		}

		switch file.Extname {
		case ".md":
			file.Extname = ".html"
			if file.Basename == "README" {
				file.Basename = "index"
			}

			// TODO: pretty path support
			// TODO: simultaneous README/index problem

			resource := whirligig.Resource{
				Title:   "Title",   // TODO: front matter
				Excerpt: "Excerpt", // TODO: front matter
				Url:     "url",     // TODO: front matter
				File:    file,
			}

			site.Resources = append(site.Resources, &resource)

		default:
			site.Files = append(site.Files, &file)
		}

		return nil
	})

	return err
}

// Returns the created timestamp of the file, if available. The modified time will be used as a fallback.
func (b *Builder) getCreatedTime(info os.FileInfo) time.Time {
	result := info.ModTime()

	// Should approximate creation time on most *nix systems
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		changeTime := time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec)

		if changeTime.Before(info.ModTime()) {
			result = changeTime
		}
	}

	return result
}
