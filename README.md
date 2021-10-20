# Gimlet Stack

Bootstrap curated Kubernetes stacks.

Logging, metrics, ingress and more - delivered with gitops.


- You can install logging aggregators, metric collectors, ingress controllers and more on your cluster with a few commands, without much knowledge of Helm charts, and their configuration options

- The components are delivered through a plain git repository with self-contained gitops automation

- You will get constant upgrades for the installed components from the stack curators

> Gimlet Stack is an open-source relaunch of the 1clickinfra.com service

## Installation

```
curl -L https://github.com/gimlet-io/gimlet-stack/releases/download/v0.3.5/stack-$(uname)-$(uname -m) -o stack
chmod +x stack
sudo mv ./stack /usr/local/bin/stack
stack --version
```

## Documentation

https://gimlet.io/gimlet-stack/getting-started/

This repo is only the tooling that delivers stacks.
If you are interested of the source of our reference stack, visit https://github.com/gimlet-io/gimlet-stack-reference

## Community

Gimlet Stack is developed in the open.

Please check the [v0.3.5](https://github.com/gimlet-io/gimlet-stack/projects/2) project and see if you are interested to help.

You can talk to us at our Discord server: https://discord.gg/VdQrdqwReB

Hope to see you there!
