# Copyright 2023 Dimitri Koshkin. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

.PHONY: mockgen
mockgen: ## Generates mock implementations
mockgen: install-tool.go.mockgen
	mockgen -source=./pkg/assistant/assistant.go -destination=./pkg/assistant/mock.go -package=assistant -copyright_file=./header.txt
	mockgen -source=./pkg/exec/exec.go -destination=./pkg/exec/mock.go -package=exec -copyright_file=./header.txt
