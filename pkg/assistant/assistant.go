// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package assistant

import (
	"context"
	"fmt"
	"strings"

	"github.com/PullRequestInc/go-gpt3"

	"github.com/dkoshkin/kubectl-assistant/pkg/exec"
)

const (
	//nolint:lll // Long line is fine here
	systemPrompt = "You should return exec commands and minimal explanation, do not include sample output. All code snippets should be surrounded by ```"
)

type Runner struct {
	gpt3Client gpt3.Client
	execRunner exec.Runner

	lastResponse string
}

func New(gpt3APIKey string, runner exec.Runner) Runner {
	return Runner{
		gpt3Client: gpt3.NewClient(gpt3APIKey),
		execRunner: runner,
	}
}

func (a *Runner) GetResponse(ctx context.Context, prompt string) (string, error) {
	resp, err := a.gpt3Client.ChatCompletion(ctx, gpt3.ChatCompletionRequest{
		Messages: []gpt3.ChatCompletionRequestMessage{
			{
				Role:    "system",
				Content: systemPrompt,
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", err
	}

	a.lastResponse = resp.Choices[0].Message.Content

	return a.lastResponse, nil
}

// CanRunKubectlCommand checks if the last response had a exec command.
func (a *Runner) CanRunKubectlCommand() bool {
	// ignore the error we only care about the command
	command, _ := a.getKubectlCommand()
	return len(command) > 0
}

// RunKubectlCommand will run the exec command from the previous response.
func (a *Runner) RunKubectlCommand() error {
	command, err := a.getKubectlCommand()
	if err != nil {
		return err
	}
	return a.execRunner.Run(command)
}

// getKubectlCommand will return the kubectl command from a code snippet (wrapped with ```).
func (a *Runner) getKubectlCommand() (string, error) {
	command, err := findCodeSnippet(a.lastResponse)
	if err != nil {
		return "", fmt.Errorf("could not parse last response for a exec command: %w", err)
	}

	if !a.execRunner.IsKubectlCommand(command) {
		return "", fmt.Errorf("command doesn't appear to be a exec command")
	}

	return command, nil
}

func findCodeSnippet(response string) (string, error) {
	start := strings.Index(response, "```")
	if start == -1 {
		return "", fmt.Errorf("did not find opening ``` in previous output")
	}
	start += 3

	end := strings.Index(response[start:], "```")
	if end == -1 {
		return "", fmt.Errorf("did not find closing ``` in previous output")
	}
	end--

	command := strings.TrimSpace(response[start : start+end])

	return command, nil
}
