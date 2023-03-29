# Copyright 2023 Dimitri Koshkin. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

.PHONY: actionlint
actionlint: ## Runs actionlint to lint Github Actions
actionlint: install-tool.go.actionlint; $(info $(M) running pre-commit)
	actionlint
