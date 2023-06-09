// Copyright 2023 Dimitri Koshkin. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkoshkin/kubectl-assistant/pkg/assistant"
	"github.com/dkoshkin/kubectl-assistant/pkg/exec"
	"github.com/dkoshkin/kubectl-assistant/pkg/prompter"
	"github.com/dkoshkin/kubectl-assistant/pkg/version"
)

const (
	apiKeyEnv = "OPENAI_KEY"
)

func main() {
	if len(os.Args) > 1 &&
		(os.Args[1] == "version" || os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Println(version.Print())
		os.Exit(0)
	}

	executor := exec.New(io.Discard, os.Stderr)
	if err := executor.KubernetesClusterRunning(); err != nil {
		log.Fatalln("Is a Kubernetes cluster running?")
	}

	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		log.Fatalln(
			"OPENAI_KEY with OpenAI key is not set, visit https://platform.openai.com/account/api-keys to create one.",
		)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	// Wait for exit signal
	go func() {
		<-sigs
		fmt.Println("Exiting...")
		os.Exit(0)
	}()

	executor = exec.New(os.Stdout, os.Stderr)
	peggy := assistant.New(apiKey, executor)
	ron := prompter.New(peggy, executor, os.Stdin, os.Stdout)

	if err := ron.Loop(); err != nil {
		log.Fatalln(err)
	}
}
