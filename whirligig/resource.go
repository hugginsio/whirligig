// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package whirligig

// Resource represents a File that will be processed by Whirligig's renderers.
type Resource struct {
	Title   string
	Excerpt string
	Url     string
	File
	// TODO: front matter data
}

func (r *Resource) Write(sourceRoot string, destinationRoot string, content []byte) error {
	destinationFile, err := r.prepareDestinationFile(destinationRoot)
	if err != nil {
		return err
	}

	defer destinationFile.Close()

	if _, err := destinationFile.Write(content); err != nil {
		return err
	}

	return nil
}
