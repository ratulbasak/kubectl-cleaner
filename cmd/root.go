package cmd

import (
	"github.com/ratulbasak/kubectl-cleaner/internal/kube"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	namespace string
	dryRun    bool
	rulesFile string

	deploymentReplicas    int
	deploymentReplicasSet bool

	jobCompletedOnly    bool
	jobCompletedOnlySet bool

	pvcPhasesRaw   string
	secretTypesRaw string

	orphanedServicesOnly    bool
	orphanedServicesOnlySet bool

	orphanedSecretsOnly    bool
	orphanedSecretsOnlySet bool

	orphanedConfigMapsOnly    bool
	orphanedConfigMapsOnlySet bool

	olderThan    int
	olderThanSet bool
)

var rootCmd = &cobra.Command{
	Use:   "kubectl-cleaner",
	Short: "Safe, intelligent cleanup of unused Kubernetes resources",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", true, "Simulate actions only (default true)")
	rootCmd.PersistentFlags().StringVar(&rulesFile, "rules-file", "", "YAML config file for analyzer rules")
	rootCmd.PersistentFlags().IntVar(&deploymentReplicas, "deployments-replicas", -1, "Consider Deployment stale if replicas <= N (-1 disables this rule)")
	rootCmd.PersistentFlags().BoolVar(&jobCompletedOnly, "jobs-completed-only", false, "Only include completed Jobs as stale")
	rootCmd.PersistentFlags().StringVar(&pvcPhasesRaw, "pvc-phases", "", "Comma-separated PVC phases to consider unused")
	rootCmd.PersistentFlags().StringVar(&secretTypesRaw, "secret-types", "", "Comma-separated secret types to consider")
	rootCmd.PersistentFlags().BoolVar(&orphanedServicesOnly, "orphaned-services-only", false, "Include only orphaned services")
	rootCmd.PersistentFlags().BoolVar(&orphanedSecretsOnly, "orphaned-secrets-only", false, "Include only orphaned secrets")
	rootCmd.PersistentFlags().BoolVar(&orphanedConfigMapsOnly, "orphaned-configmaps-only", false, "Include only orphaned configmaps")
	rootCmd.PersistentFlags().IntVar(&olderThan, "older-than", -1, "Mark resources as stale if older than N days (-1 disables this rule)")
	// Record which flags user set
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		deploymentReplicasSet = cmd.Flags().Changed("deployments-replicas")
		jobCompletedOnlySet = cmd.Flags().Changed("jobs-completed-only")
		orphanedServicesOnlySet = cmd.Flags().Changed("orphaned-services-only")
		orphanedSecretsOnlySet = cmd.Flags().Changed("orphaned-secrets-only")
		orphanedConfigMapsOnlySet = cmd.Flags().Changed("orphaned-configmaps-only")
		olderThanSet = cmd.Flags().Changed("older-than")
	}
}

func prepareRules(cmd *cobra.Command) kube.AnalyzerRules {
	rules := kube.DefaultAnalyzerRules()
	if rulesFile != "" {
		f, err := os.Open(rulesFile)
		if err == nil {
			rules = kube.LoadAnalyzerRules(f)
			f.Close()
		}
	}
	if deploymentReplicasSet {
		rules.DeploymentReplicas = &deploymentReplicas
	}
	if jobCompletedOnlySet {
		rules.JobCompletedOnly = &jobCompletedOnly
	}
	if pvcPhasesRaw != "" {
		rules.PVCPhases = parseCommaList(pvcPhasesRaw)
	}
	if secretTypesRaw != "" {
		rules.SecretTypes = parseCommaList(secretTypesRaw)
	}
	if orphanedServicesOnlySet {
		rules.OrphanedServicesOnly = &orphanedServicesOnly
	}
	if orphanedSecretsOnlySet {
		rules.OrphanedSecretsOnly = &orphanedSecretsOnly
	}
	if orphanedConfigMapsOnlySet {
		rules.OrphanedConfigMapsOnly = &orphanedConfigMapsOnly
	}
	if olderThanSet {
		rules.OlderThanDays = &olderThan
	}
	return rules
}

func parseCommaList(s string) []string {
	out := []string{}
	for _, v := range strings.Split(s, ",") {
		if trimmed := strings.TrimSpace(v); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
