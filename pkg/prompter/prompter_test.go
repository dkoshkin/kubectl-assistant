// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package prompter

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dkoshkin/kubectl-assistant/pkg/assistant"
	"github.com/dkoshkin/kubectl-assistant/pkg/exec"
)

func TestPrompter_Loop_runKubectlDirectly(t *testing.T) {
	ctrl := gomock.NewController(t)

	in := strings.NewReader(`
kubectl get nodes
`)

	mockAssistantRunner := assistant.NewMockRunner(ctrl)
	mockExecRunner := exec.NewMockRunner(ctrl)

	mockExecRunner.EXPECT().IsKubectlCommand("").Return(false).Times(1)
	mockExecRunner.EXPECT().IsKubectlCommand("kubectl get nodes").Return(true).Times(1)
	mockExecRunner.EXPECT().Run("kubectl get nodes").Return(nil).Times(1)

	out := bytes.Buffer{}
	expectedOutput := []byte(
		//nolint:lll // Long lines are fine in tests
		"Begin by typing what you want to accomplish in your Kubernetes cluster and then hit \"Enter\".\nFor example:\n  List all control-plane Nodes\n  Get Kubernetes versions for all Nodes\n  Create deployment named nginx, using image nginx and ports 80 and 443\n  Find all objects with label app=nginx\n\nYou will then see some text output and in most cases either a exec command or some YAML output.\nIf the command looks reasonable to you, type in \"k\" and then hit \"Enter\" to execute it against the cluster.\n\nYou can also type \"kubectl ...\" to execute a custom command.\n\nHit CTRL+C to exit.\n\n> \n> \n> ",
	)

	p := New(mockAssistantRunner, mockExecRunner, in, &out)
	err := p.Loop()
	assert.EqualError(t, err, "an error occurred reading input: EOF", "expected an EOF error")
	assert.Equal(t, string(expectedOutput), out.String())
}

func TestPrompter_Loop_runPromptOutput(t *testing.T) {
	ctrl := gomock.NewController(t)

	in := strings.NewReader(`
list all Nodes
k
`)

	mockAssistantRunner := assistant.NewMockRunner(ctrl)
	mockExecRunner := exec.NewMockRunner(ctrl)

	// expect a call to exec.Run("kubectl get nodes") in assistant.RunKubectlCommand() is called
	mockAssistantRunner.EXPECT().RunKubectlCommand().DoAndReturn(func() error {
		//nolint:wrapcheck // Want to return the error as is in the mock.
		return mockExecRunner.Run("kubectl get nodes")
	})
	mockAssistantRunner.EXPECT().
		GetResponse(context.TODO(), "list all Nodes").
		Return("To list all nodes, run the command:\n\n```\nkubectl get nodes\n```", nil).
		Times(1)
	mockAssistantRunner.EXPECT().CanRunKubectlCommand().Return(true).Times(1)

	mockExecRunner.EXPECT().IsKubectlCommand("").Return(false).Times(1)
	mockExecRunner.EXPECT().IsKubectlCommand("list all Nodes").Return(false).Times(1)
	mockExecRunner.EXPECT().IsKubectlCommand("k").Return(false).Times(1)
	mockExecRunner.EXPECT().Run("kubectl get nodes").Return(nil).Times(1)

	out := bytes.Buffer{}
	expectedOutput := []byte(
		//nolint:lll // Long lines are fine in tests
		"Begin by typing what you want to accomplish in your Kubernetes cluster and then hit \"Enter\".\nFor example:\n  List all control-plane Nodes\n  Get Kubernetes versions for all Nodes\n  Create deployment named nginx, using image nginx and ports 80 and 443\n  Find all objects with label app=nginx\n\nYou will then see some text output and in most cases either a exec command or some YAML output.\nIf the command looks reasonable to you, type in \"k\" and then hit \"Enter\" to execute it against the cluster.\n\nYou can also type \"kubectl ...\" to execute a custom command.\n\nHit CTRL+C to exit.\n\n> \n> ====================================================================================================\nTo list all nodes, run the command:\n\n```\nkubectl get nodes\n```\n\n> \n> ",
	)

	p := New(mockAssistantRunner, mockExecRunner, in, &out)
	err := p.Loop()
	assert.EqualError(t, err, "an error occurred reading input: EOF", "expected an EOF error")
	assert.Equal(t, string(expectedOutput), out.String())
}
