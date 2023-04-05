// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"io"
	"os/exec"
	"strings"
)

type Runner interface {
	KubernetesClusterRunning() error
	Run(command string) error
	IsKubectlCommand(in string) bool
}

type BashRunner struct {
	out io.Writer
	err io.Writer
}

func New(out, err io.Writer) Runner {
	return BashRunner{out: out, err: err}
}

func (r BashRunner) KubernetesClusterRunning() error {
	return r.Run("kubectl cluster-info")
}

// Run will run the command with 'bash -c'.
func (r BashRunner) Run(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = r.out
	cmd.Stderr = r.err
	//nolint:wrapcheck // Want to return the error as is.
	return cmd.Run()
}

// IsKubectlCommand returns true if in starts with "kubectl "
// TODO: make this check more robust.
func (r BashRunner) IsKubectlCommand(in string) bool {
	return strings.HasPrefix(in, "kubectl ")
}
