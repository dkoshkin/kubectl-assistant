// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package assistant

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/PullRequestInc/go-gpt3"

	"github.com/dkoshkin/kubectl-assistant/pkg/exec"
)

const (
	//nolint:lll // Long line is fine here
	systemPrompt = "You should return kubectl commands and minimal explanation, do not include sample output. All code snippets should be surrounded by ```"
)

type Runner interface {
	GetResponse(ctx context.Context, prompt string) (string, error)
	CanRunKubectlCommand() bool
	RunKubectlCommand() error
}

type GPT3Runner struct {
	gpt3Client gpt3.Client
	execRunner exec.Runner

	lastResponse string
}

func New(gpt3APIKey string, runner exec.Runner) Runner {
	return &GPT3Runner{
		gpt3Client: gpt3.NewClient(gpt3APIKey),
		execRunner: runner,
	}
}

func (a *GPT3Runner) GetResponse(ctx context.Context, prompt string) (string, error) {
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
		return "", fmt.Errorf("could not get response from OpenAI: %w", err)
	}

	a.lastResponse = resp.Choices[0].Message.Content

	return a.lastResponse, nil
}

// CanRunKubectlCommand checks if the last response had a exec command.
func (a *GPT3Runner) CanRunKubectlCommand() bool {
	// ignore the error we only care about the command
	command, _ := a.getKubectlCommand()
	return len(command) > 0
}

// RunKubectlCommand will run the exec command from the previous response.
func (a *GPT3Runner) RunKubectlCommand() error {
	command, err := a.getKubectlCommand()
	if err != nil {
		return err
	}
	//nolint:wrapcheck // Want to return the error as is.
	return a.execRunner.Run(command)
}

// getKubectlCommand will return the kubectl command from a code snippet (wrapped with ```).
func (a *GPT3Runner) getKubectlCommand() (string, error) {
	command, err := findCodeSnippet(a.lastResponse)
	if err != nil {
		return "", fmt.Errorf("could not parse last response for a exec command: %w", err)
	}

	if !a.execRunner.IsKubectlCommand(command) {
		//nolint:goerr113 // No need to return a custom error.
		return "", errors.New("command doesn't appear to be a exec command")
	}

	return command, nil
}

func findCodeSnippet(response string) (string, error) {
	start, err := findSnippetOpening(response)
	if err != nil {
		return "", err
	}
	end, err := findSnippetClosing(response[start:])
	if err != nil {
		return "", err
	}

	command := strings.TrimSpace(response[start : start+end])

	return command, nil
}

func findSnippetOpening(str string) (int, error) {
	validStarts := []string{"```bash", "```shell", "```"}

	start := -1
	for _, substr := range validStarts {
		start = strings.Index(str, substr)
		if start != -1 {
			start += len(substr)
			break
		}
	}
	if start == -1 {
		//nolint:goerr113 // No need to return a custom error.
		return start, fmt.Errorf(
			"did not find opening string in previous output, looking for: %s",
			strings.Join(validStarts, ", "),
		)
	}

	return start, nil
}

func findSnippetClosing(str string) (int, error) {
	end := strings.Index(str, "```")
	//nolint:goerr113 // No need to return a custom error.
	if end == -1 {
		return end, fmt.Errorf("did not find closing ``` in previous output")
	}
	// exclude a newline character
	if str[end] == '\n' {
		end--
	}

	return end, nil
}
