// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package assistant

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_findCodeSnippet(t *testing.T) {
	tests := []struct {
		name            string
		response        string
		expectedCommand string
		expectedErr     error
	}{
		{
			name:            "response with command",
			response:        "To list all nodes, run the command:\n\n```\nkubectl get nodes\n```",
			expectedCommand: "kubectl get nodes",
		},
		{
			name:            "response with command wrapped in ```bash",
			response:        "To list all nodes, run the command:\n\n```bash\nkubectl get nodes\n```",
			expectedCommand: "kubectl get nodes",
		},
		{
			name: "response with command wrapped in ``` without a newline",
			//nolint:lll // Long lines are fine in tests
			response:        "To show deployments with less than desired pods ready, you can use the following command:\n\n```kubectl get deployments --field-selector=\"status.replicas!=status.readyReplicas\"```",
			expectedCommand: "kubectl get deployments --field-selector=\"status.replicas!=status.readyReplicas\"",
		},
		{
			name: "response without a command",
			//nolint:lll // Long lines are fine in tests
			response: "Kubernetes is an open-source container orchestration platform that automates the deployment, scaling, and management of containerized applications. It enables operators to manage applications across multiple nodes, and provides mechanisms for the automated deployment, scaling, and recovery of instances of those applications.",
			expectedErr: fmt.Errorf(
				"did not find opening string in previous output, looking for: ```bash, ```shell, ```",
			),
		},
	}
	for _, tt := range tests {
		gotCommand, gotErr := findCodeSnippet(tt.response)
		assert.Equal(t, tt.expectedErr, gotErr)
		assert.Equal(t, tt.expectedCommand, gotCommand)
	}
}
