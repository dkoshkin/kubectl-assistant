// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package version

import (
	"bytes"
	"html/template"
	"runtime"
	"strings"
)

// Build information. Populated at build-time.
//
//nolint:gochecknoglobals // Version globals are set at build time.
var (
	version      string
	major        string
	minor        string
	patch        string
	revision     string
	branch       string
	commitDate   string
	gitTreeState string
	goVersion    = runtime.Version()
)

// Print returns version information.
func Print() string {
	m := map[string]string{
		"version":      version,
		"major":        major,
		"minor":        minor,
		"patch":        patch,
		"revision":     revision,
		"branch":       branch,
		"commitDate":   commitDate,
		"gitTreeState": gitTreeState,
		"goVersion":    goVersion,
		"platform":     runtime.GOOS + "/" + runtime.GOARCH,
	}
	t := template.Must(template.New("version").Parse(`
	version {{.version}} (branch: {{.branch}}, revision: {{.revision}}{{with .gitTreeState}}, gitTreeState: {{.}}{{end}})
		build date:       {{.commitDate}}
		go version:       {{.goVersion}}
		platform:         {{.platform}}
	`))

	var buf bytes.Buffer
	if err := t.ExecuteTemplate(&buf, "version", m); err != nil {
		panic(err)
	}
	return strings.TrimSpace(buf.String())
}
