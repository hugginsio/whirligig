// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package whirligig

import (
	"time"
)

type Site struct {
	Time      time.Time   // The current time when the site is built.
	Resources []*Resource // A list of all Resources in the source directory. Resources are processed by Whirligig's renderers.
	Files     []*File     // A list of all Files in the source directory. Files are copied as-is and not processed.
	Version   string      // The commit hash of the site, if it is in a Git repository.
	// TODO: data
	// TODO: url
}
