# Gimlet Stack

Bootstrap curated Kubernetes stacks.

Logging, metrics, ingress and more - delivered with gitops.

> Gimlet Stack is an open-source relaunch of the 1clickinfra.com service

## Goals

- The common infrastructure elements should be provisioned on any Kubernetes cluster without having to set too many flags
- The components are integrated. Eg.: one Grafana configured to monitor all components
- Sane defaults that are tailored for small to mid-sized teams: 1-50 developers
- Cloud provider flavors. Ingress annotations, storage classes differ from cloud-to-cloud. The goal is to cover those differences
- GUI
- Upgrade paths
- Multiple curators: unlike 1clickinfra.com, Gimlet Stack should support many curators. You can pick who you trust and take their stack.
This repo contains the delivery engine, not the stacks.

## Community

Gimlet Stack is developed in the open.

Please check the [v0.1.0](https://github.com/gimlet-io/gimlet-stack/projects/1) project and see if you are interested to help.

You can talk to us at our Discord server: https://discord.gg/VdQrdqwReB

Hope to see you there!
