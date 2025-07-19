# kubectl-cleaner

A safe, interactive, and intelligent Krew plugin to analyze and clean up unused Kubernetes resources (Deployments, PVCs, Jobs, Secrets, ConfigMaps) in a namespace.

[![Go](https://img.shields.io/badge/Go-1.22+-brightgreen)](https://golang.org)
[![Krew Plugin](https://img.shields.io/badge/Krew-Plugin-blueviolet)](https://krew.sigs.k8s.io/docs/)

---

## Features

- **Analyze** or **purge**: Dry-run analysis or actual cleanup of resources.
- **Interactive rules**: Control what’s considered “stale” or “orphaned” via CLI flags or YAML config.
- **Resources covered**:
    - Deployments (by replicas/age)
    - StatefulSets, DaemonSets, Jobs, CronJobs, PVCs (by age)
    - Services (finds orphaned)
    - Orphaned ConfigMaps & Secrets (with reference detection)
- **Safe defaults** (dry-run, conservative filters)
- **Works great as a [kubectl krew plugin](https://krew.sigs.k8s.io/docs/)**
- **Mac (arm64, amd64) & Linux (amd64) builds**

---

## Installation

```shell
kubectl krew install cleaner
```
Or, to test your local build:
```shell
kubectl krew install --manifest=plugin.yaml --archive=kubectl-cleaner_darwin_arm64.tar.gz
```

## Usage

Analyze your namespace for stale or orphaned resources:
```shell
kubectl cleaner analyze --namespace=my-ns --older-than=14
```

Purge (delete) unused resources:
```shell
kubectl cleaner purge --namespace=dev --older-than=14 --dry-run=false
```

Use a config file for custom rules:
```shell
kubectl cleaner analyze --namespace=prod --rules-file=rules.yaml
```


## CLI Flags
| Flag                         | Description                                                       |
| ---------------------------- | ----------------------------------------------------------------- |
| `--namespace`                | Namespace to analyze (default: default)                           |
| `--older-than`               | Mark resources as stale if older than N days (optional)           |
| `--deployments-replicas`     | Mark Deployment as stale if replicas <= N (-1 disables this rule) |
| `--jobs-completed-only`      | Only include completed Jobs as stale                              |
| `--pvc-phases`               | Comma-separated PVC phases (e.g. Released,Lost)                   |
| `--secret-types`             | Comma-separated secret types (e.g. Opaque)                        |
| `--orphaned-secrets-only`    | Only include orphaned secrets                                     |
| `--orphaned-configmaps-only` | Only include orphaned configmaps                                  |
| `--orphaned-services-only`   | Only include orphaned services                                    |
| `--dry-run`                  | If true, only simulates (default: true)                           |
| `--rules-file`               | YAML config file for advanced rule config                         |



Example `rules.yaml`
```yaml
deployment_replicas: 0
job_completed_only: true
pvc_phases: ["Released", "Lost"]
secret_types: ["Opaque"]
orphaned_secrets_only: true
orphaned_configmaps_only: true
orphaned_services_only: true
older_than: 14
```

Help Command

```shell
kubectl cleaner -h
Safe, intelligent cleanup of unused Kubernetes resources

Usage:
  kubectl-cleaner [command]

Available Commands:
  analyze     Analyze namespace for stale or orphaned resources
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  purge       Delete unused/orphaned resources

Flags:
      --deployments-replicas int   Consider Deployment stale if replicas <= N (-1 disables this rule) (default -1)
      --dry-run                    Simulate actions only (default true) (default true)
  -h, --help                       help for kubectl-cleaner
      --jobs-completed-only        Only include completed Jobs as stale
      --namespace string           Kubernetes namespace (default "default")
      --older-than int             Mark resources as stale if older than N days (-1 disables this rule) (default -1)
      --orphaned-configmaps-only   Include only orphaned configmaps
      --orphaned-secrets-only      Include only orphaned secrets
      --orphaned-services-only     Include only orphaned services
      --pvc-phases string          Comma-separated PVC phases to consider unused
      --rules-file string          YAML config file for analyzer rules
      --secret-types string        Comma-separated secret types to consider

Use "kubectl-cleaner [command] --help" for more information about a command.
```
