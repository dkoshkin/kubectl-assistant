<!--
 Copyright 2023 Dimitri Koshkin. All rights reserved.
 SPDX-License-Identifier: Apache-2.0
 -->

# kubectl-assistant

This tool uses [OpenAI's GPT3](https://platform.openai.com/account/api-keys) API to generate Kubernetes [kubectl](https://kubernetes.io/docs/reference/kubectl/) commands to run against a cluster.
It aims to help find `kubectl` commands for what you are trying to accomplish, without needing to search around in the docs.

## Prerequisites

- A running Kubernetes cluster
- [kubectl](https://kubernetes.io/docs/reference/kubectl/) in your `PATH`

## Usage Instructions

1.  Download the binary for your OS and architecture.

2.  Move it somewhere in your `PATH`.

3.  Generate an OpenAI API Key [here](https://platform.openai.com/account/api-keys).

4.  Run `kubectl assistant`

    ```txt
    $ export OPENAI_KEY=<>
    $ kubectl assistant
    Begin by typing what you want to accomplish in your Kubernetes cluster and then hit "Enter".
    For example:
      List all control-plane Nodes
      Get Kubernetes versions for all Nodes
      Create deployment named nginx, using image nginx and ports 80 and 443
      Find all objects with label app=nginx

    You will then see some text output and in most cases either a exec command or some YAML output.
    If the command looks reasonable to you, type in "k" and then hit "Enter" to execute it against the cluster.

    You can also type "kubectl ..." to execute a custom command.

    Hit CTRL+C to exit.

    > List all control-plane Nodes
    ==============================================================================================================================================================================================================
    To list all control-plane nodes in your Kubernetes cluster, you can use the following exec command:

    ```bash
    kubectl get nodes --selector=node-role.kubernetes.io/control-plane
    ``

    This command will display a list of all the control-plane nodes in the cluster. The `--selector=node-role.kubernetes.io/control-plane` flag filters the list of nodes based on a label selector that is automatically applied to control-plane nodes by default when they are registered with the cluster.

    > k
    NAME                 STATUS   ROLES           AGE     VERSION
    kind-control-plane   Ready    control-plane   3d16h   v1.25.3

    > kubectl get pods -A
    NAMESPACE            NAME                                         READY   STATUS    RESTARTS   AGE
    default              nginx                                        1/1     Running   0          3d16h
    default              nginx-56fd7f4d49-v67dd                       1/1     Running   0          3d15h
    kube-system          coredns-565d847f94-kngc2                     1/1     Running   0          3d16h
    kube-system          coredns-565d847f94-ngfzd                     1/1     Running   0          3d16h
    kube-system          etcd-kind-control-plane                      1/1     Running   0          3d16h
    kube-system          kindnet-t69vm                                1/1     Running   0          3d16h
    kube-system          kube-apiserver-kind-control-plane            1/1     Running   0          3d16h
    kube-system          kube-controller-manager-kind-control-plane   1/1     Running   0          3d16h
    kube-system          kube-proxy-6mb6r                             1/1     Running   0          3d16h
    kube-system          kube-scheduler-kind-control-plane            1/1     Running   0          3d16h
    local-path-storage   local-path-provisioner-684f458cdd-tscq6      1/1     Running   0          3d16h
    ```

After typing in a question in the prompt you will see a response, many of them having example `kubectl ...` commands.

-   Type in `k` to execute the first command against your Kubernetes cluster.

-   Type in another question to get additional responses. (NOTE, unlike ChatGPT the conversation threads are not contextual, if you want more detail for a previois response, you must re-ask the question).

-   Type in `kubectl ...` commands directly to have them execute against the cluster.

## Setup your Dev Environment

- Install [asdf](https://github.com/asdf-community/asdf-direnv)
- Install [asdf-direnv](https://github.com/asdf-community/asdf-direnv#setup)
- Add a global `direnv` version with: `asdf global direnv latest`
- Install all tools with: `make install-tools`

Tip: to see all available make targets with descriptions, simply run `make`.

### Lint

```bash
make lint
```

### Test

```bash
make tst
```

### Build

The binary for your OS will be placed in `./dist`, e.g. `./dist/kubectl-assistant_darwin_arm64/kubectl-assisant`:

```bash
make build-snapshot
```

### Pre-commit

```bash
make pre-commit
```
