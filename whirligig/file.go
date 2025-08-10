// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package whirligig

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

// File represents all files in the source directory.
type File struct {
	Name     string         // The name of the File, e.g. `styles.css`.
	Basename string         // The basename of the File, e.g. `styles`.
	Extname  string         // The extension of the File, e.g. `.css`.
	Created  time.Time      // When the File was created.
	Modified time.Time      // When the File was modified.
	Path     string         // The path to the raw File, relative to the source directory.
	Data     map[string]any // Data extracted from companion YAML files or frontmatter.
}

func (f *File) prepareDestinationFile(destinationRoot string) (*os.File, error) {
	destFilePath := filepath.Join(destinationRoot, f.Path, f.Basename+f.Extname)

	if err := os.MkdirAll(filepath.Dir(destFilePath), 0755); err != nil {
		return nil, err
	}

	return os.Create(destFilePath)
}

// Provided the sourceRoot and destinationRoot, Copy will use the File metadata to copy the file.
func (f *File) Copy(sourceRoot string, destinationRoot string) error {
	sourceFile, err := os.Open(filepath.Join(sourceRoot, f.Path, f.Name))
	if err != nil {
		return err
	}

	defer sourceFile.Close()

	destinationFile, err := f.prepareDestinationFile(destinationRoot)
	if err != nil {
		return err
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// Content will return the contents of the File as bytes.
func (f *File) Content(sourceRoot string) ([]byte, error) {
	return os.ReadFile(filepath.Join(sourceRoot, f.Path, f.Name))
}
