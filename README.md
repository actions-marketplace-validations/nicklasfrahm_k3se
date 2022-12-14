# k3se 🧀

[![Go Report Card](https://goreportcard.com/badge/github.com/nicklasfrahm/k3se?style=flat-square)](https://goreportcard.com/report/github.com/nicklasfrahm/k3se)
[![Release](https://img.shields.io/github/release/nicklasfrahm/k3se.svg?style=flat-square)](https://github.com/nicklasfrahm/k3se/releases/latest)
[![Go Reference](https://img.shields.io/badge/Go-reference-informational.svg?style=flat-square)](https://pkg.go.dev/github.com/nicklasfrahm/k3se)

A lightweight Kubernetes engine that deploys `k3s` clusters declaratively based on a cluster configuration file. The name is an abbreviation for _k3s engine_ and a hommage to the German word for cheese, _Käse [ˈkɛːzə]_.

## Quickstart 💡

If you want to test `k3se` you can use [Vagrant][website-vagrant]. All examples in the `examples/` folder can be used with the provided `Vagrantfile` that provisions 3 Ubuntu VMs. To bring up the VMs you can run the following command:

```bash
make vagrant-up
```

Once you are done testing, you can destroy the VMs with the following command:

```bash
make vagrant-down
```

## Prerequisites 📝

The nodes have to be accessible via SSH, either directly or via a bastion host. Further, the user on the remote nodes needs to have passwordless `sudo` set up. If this is not yet the case, you may manually do so via the following command:

```bash
echo "$(whoami) ALL=(ALL) NOPASSWD: ALL" | sudo tee /etc/sudoers.d/$(whoami)
```

## Limitations 🚨

The following features are currently not supported, but are planned for future releases:

- **Downsizing**  
  If nodes are removed from the cluster configuration, they are not decommissioned. We plan to enable this feature in the future; performing operations similar to `kubectl cordon` and `kubectl drain` automagically.
- **Diffing**  
  Using the `git` history of the cluster configuration to display potential actions that can be taken to bring the cluster up to date with the configuration.

## GitHub Actions 🤖

You can use GitHub Actions to run `k3se` inside your CI pipeline. Simply add the following workflow:

```yaml
name: Deploy

on:
  push:
    branches:
      - main

jobs:
  kubernetes:
    name: Kubernetes
    runs-on: ubuntu-latest
    steps:
      - name: Check out repository
        uses: actions/checkout@v3

      - name: Deploy k3s cluster
        uses: nicklasfrahm/k3se@main
        with:
          command: up examples/standalone.yml
```

### Input variables ⚙️

See [action.yml](./action.yml) for more detailed information. **Please note that all input variables must have string values. It is thus recommend to always use quotes.**

| Input variable | Default value | Description                                                                                                    |
| -------------- | ------------- | -------------------------------------------------------------------------------------------------------------- |
| `command`      | `up`          | Subcommand to be executed. Use `up path/to/my/config.yaml` to specify a custom location for your cluster spec. |
| `version`      | `latest`      | Version of `k3se` to use.                                                                                      |

## License 📄

This project is and will always be licensed under the terms of the [MIT license][file-license].

[file-license]: https://www.apache.org/licenses/LICENSE-2.0
[website-vagrant]: https://vagrantup.com
