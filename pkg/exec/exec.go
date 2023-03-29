// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package exec

import (
	"io"
	"os/exec"
	"strings"
)

type Runner struct {
	out io.Writer
	err io.Writer
}

func New(out, err io.Writer) Runner {
	return Runner{out: out, err: err}
}

func (r Runner) KubernetesClusterRunning() error {
	return r.Run("kubectl cluster-info")
}

// Run will run the command with 'bash -c'.
func (r Runner) Run(command string) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = r.out
	cmd.Stderr = r.err
	return cmd.Run()
}

// IsKubectlCommand returns true if the command starts with "exec "
// TODO: make this check more robust.
func (r Runner) IsKubectlCommand(command string) bool {
	return strings.HasPrefix(command, "kubectl ")
}
