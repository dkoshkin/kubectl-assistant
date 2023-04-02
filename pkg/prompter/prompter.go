// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package prompter

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"golang.org/x/term"

	"github.com/dkoshkin/kubectl-assistant/pkg/assistant"
	"github.com/dkoshkin/kubectl-assistant/pkg/exec"
)

const (
	welcomeText = `Begin by typing what you want to accomplish in your Kubernetes cluster and then hit "Enter".
For example:
  List all control-plane Nodes
  Get Kubernetes versions for all Nodes
  Create deployment named nginx, using image nginx and ports 80 and 443
  Find all objects with label app=nginx

You will then see some text output and in most cases either a exec command or some YAML output.
If the command looks reasonable to you, type in "k" and then hit "Enter" to execute it against the cluster.

You can also type "exec ..." to execute a custom command.

Hit CTRL+C to exit.
`
)

type Prompter struct {
	assistantRunner assistant.Runner
	execRunner      exec.Runner

	in  io.Reader
	out io.Writer
}

func New(
	assistantRunner assistant.Runner,
	execRunner exec.Runner,
	in io.Reader,
	out io.Writer,
) Prompter {
	return Prompter{
		assistantRunner: assistantRunner,
		execRunner:      execRunner,
		in:              in,
		out:             out,
	}
}

func (p *Prompter) Loop() error {
	fmt.Print(welcomeText)

	ctx := context.Background()
	reader := bufio.NewReader(p.in)

	// continue until user exits
	for {
		fmt.Fprint(p.out, "\n> ")

		// ReadString will block until a newline character
		prompt, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("an error occurred reading input: %w", err)
		}
		prompt = strings.TrimSpace(strings.ReplaceAll(prompt, "\n", ""))

		// convert if else to switch
		switch {
		case p.execRunner.IsKubectlCommand(prompt): // run exec commands directly if provided by user
			err = p.execRunner.Run(prompt)
			if err != nil {
				fmt.Fprintln(
					p.out,
					"An error occurred running exec command. Please try again",
					err,
				)
			}
		case prompt == "k": // ignore if a user just hits "ENTER"
			if !p.assistantRunner.CanRunKubectlCommand() {
				fmt.Fprintln(
					p.out,
					"Previous output did not contain any commands to run. Please try again",
				)
				continue
			}
			err = p.assistantRunner.RunKubectlCommand()
			if err != nil {
				fmt.Fprintln(p.out, "An error occurred running command. Please try again", err)
				continue
			}
			continue
		case prompt == "": // ignore if a user just hits "ENTER"
			continue
		default:
			// get the response from the assistant
			resp, err := p.assistantRunner.GetResponse(ctx, prompt)
			if err != nil {
				fmt.Fprintln(p.out, "An error occurred generating output. Please try again", err)
				continue
			}

			// print a separator and then the response from the assistant
			fmt.Fprintln(p.out, strings.Repeat("=", terminalWidth()))
			fmt.Fprintln(p.out, resp)
		}
	}
}

func terminalWidth() int {
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 100
	}
	return width
}
