package cmd

import (
	"fmt"
	"github.com/ratulbasak/kubectl-cleaner/internal/kube"
	"github.com/spf13/cobra"
)

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze namespace for stale or orphaned resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := kube.GetKubeClient()
		if err != nil {
			return err
		}
		rules := prepareRules(cmd)
		report, err := kube.AnalyzeNamespace(client, namespace, rules)
		if err != nil {
			return err
		}
		if len(report) == 0 {
			fmt.Println("No unused, orphaned, or stale resources found.")
			return nil
		}
		fmt.Println("Potentially removable resources:")
		for _, res := range report {
			fmt.Printf("- %-12s %s\n", res.Kind, res.Name)
		}
		fmt.Println("\n(Dry-run only. No resources have been deleted.)")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
}
