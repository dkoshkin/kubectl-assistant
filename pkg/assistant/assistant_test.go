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
		name          string
		response      string
		expectCommand string
		expectErr     error
	}{
		{
			name:          "response with command",
			response:      "To list all nodes, run the command:\n\n```\nkubectl get nodes\n```",
			expectCommand: "kubectl get nodes",
		},
		{
			name:          "response with command wrapped in ```bash",
			response:      "To list all nodes, run the command:\n\n```bash\nkubectl get nodes\n```",
			expectCommand: "kubectl get nodes",
		},
		{
			name: "response without a command",
			//nolint:lll // Long lines are fine in tests
			response: "Kubernetes is an open-source container orchestration platform that automates the deployment, scaling, and management of containerized applications. It enables operators to manage applications across multiple nodes, and provides mechanisms for the automated deployment, scaling, and recovery of instances of those applications.",
			expectErr: fmt.Errorf(
				"did not find opening string in previous output, looking for: ```bash, ```shell, ```",
			),
		},
	}
	for _, tt := range tests {
		gotCommand, gotErr := findCodeSnippet(tt.response)
		assert.Equal(t, tt.expectErr, gotErr)
		assert.Equal(t, tt.expectCommand, gotCommand)
	}
}
