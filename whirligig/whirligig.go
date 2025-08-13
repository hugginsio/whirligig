// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package whirligig

import goversion "github.com/caarlos0/go-version"

type Whirligig struct {
	SourcePath  string          // The site source directory.
	Version     string          // The version of Whirligig used to build the site.
	VersionInfo *goversion.Info // The full Whirligig version information.
}
