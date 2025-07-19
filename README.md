# kubectl-cleaner

A safe, interactive, and intelligent Krew plugin to analyze and clean up unused Kubernetes resources (Deployments, PVCs, Jobs, Secrets, ConfigMaps) in a namespace.

## Usage

```shell
kubectl cleaner --namespace=my-ns --dry-run
kubectl cleaner --namespace=prod-ns --dry-run=false
```

Help

```shell
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