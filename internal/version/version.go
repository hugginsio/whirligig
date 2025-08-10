// Copyright (c) Kyle Huggins
// SPDX-License-Identifier: BSD-3-Clause

package version

import goversion "github.com/caarlos0/go-version"

func GetVersionInfo() goversion.Info {
	return goversion.GetVersionInfo()
}
