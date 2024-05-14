# ynv

ynv is a utilizty that allows you to easily large YAML files containing multiple documents. It's primary use-case
is when working with [Helm](https://helm.sh/) or other Kubernetes tools that output multiple YAML documents.

## Installation

### MacOS

```shell
brew install snooc/ynv/ynv
```

## Usage

```
$ helm template my-release my-chart --values my-values.yaml | ynv
postgresql/templates/primary/networkpolicy.yaml
postgresql/templates/serviceaccount.yaml
postgresql/templates/secrets.yaml
postgresql/templates/primary/svc-headless.yaml
postgresql/templates/primary/svc.yaml
postgresql/templates/primary/statefulset.yaml
```

### FZF Support

ynv also support interactive mode if you have [fzf](https://github.com/junegunn/fzf) installed. Interactive mode can be enabled using the `-i` flag or setting the environment variable `YNV_INTERACTIVE=true`.

```
$ helm template my-release my-chart --values my-values.yaml | ynv -i
```
