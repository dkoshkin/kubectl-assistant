# Copyright 2023 Dimitri Koshkin. All rights reserved.
# SPDX-License-Identifier: Apache-2.0

INCLUDE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

include $(INCLUDE_DIR)github-actions.mk
include $(INCLUDE_DIR)go.mk
include $(INCLUDE_DIR)goreleaser.mk
include $(INCLUDE_DIR)help.mk
include $(INCLUDE_DIR)make.mk
include $(INCLUDE_DIR)pre-commit.mk
include $(INCLUDE_DIR)repo.mk
include $(INCLUDE_DIR)tools.mk
include $(INCLUDE_DIR)upx.mk
