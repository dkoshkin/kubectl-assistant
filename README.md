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
- [asdf](https://asdf-vm.com/) and make for building locally.

## Usage Instructions

1.  Download the binary for your OS and architecture.

2.  Move it somewhere in your `PATH`.

3.  Generate an API Key [here](https://platform.openai.com/account/api-keys).

4.  Run `kubectl assistant`

    ```shell
    $ exec assistant
    Begin by typing what you want to accomplish in your Kubernetes cluster and then hit "Enter".
    For example:
      List all control-plane Nodes
      List all Pods that don't have an ImagePullPolicy of Always
      Create deployment named nginx, using image nginx and ports 80 and 443

    You will then see some text output and in most cases either a exec command or some YAML output.
    If the command looks reasonable to you, type in "k" and then hit "Enter" to execute it against the cluster.

    You can also type "exec ..." to execute a custom command.

    Hit CTRL+C to exit.

    > List all control-plane Nodes
    ==============================================================================================================================================================================================================
    To list all control-plane nodes in your Kubernetes cluster, you can use the following exec command:

    ```&nbsp;
    exec get nodes --selector=node-role.kubernetes.io/control-plane
    ```&nbsp;

    This command will display a list of all the control-plane nodes in the cluster. The `--selector=node-role.kubernetes.io/control-plane` flag filters the list of nodes based on a label selector that is automatically applied to control-plane nodes by default when they are registered with the cluster.

    > k
    NAME                 STATUS   ROLES           AGE     VERSION
    kind-control-plane   Ready    control-plane   3d16h   v1.25.3

    > exec get pods -A
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

## Dev Instructions

-   Built it, the binary for your OS will be placed in `./dist`, e.g. `./dist/kubectl-assistant_darwin_arm64/kubectl-assisant`:

    ```shell
    make build-snapshot
    ```

-   Test it:

    ```shell
    make test
    ```

-   Lint it:

    ```shell
    make lint
    ```
