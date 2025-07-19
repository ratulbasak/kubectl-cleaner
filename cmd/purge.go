package cmd

import (
	"fmt"
	"github.com/ratulbasak/kubectl-cleaner/internal/kube"
	"github.com/spf13/cobra"
)

var purgeCmd = &cobra.Command{
	Use:   "purge",
	Short: "Delete unused/orphaned resources",
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
		fmt.Printf("Deleting %d resources:\n", len(report))
		for _, res := range report {
			fmt.Printf("- %-12s %s\n", res.Kind, res.Name)
			if !dryRun {
				err := kube.DeleteResource(client, namespace, res)
				if err != nil {
					fmt.Printf("  ERROR: %v\n", err)
				}
			}
		}
		if dryRun {
			fmt.Println("\nDry-run mode: no resources were actually deleted.")
		} else {
			fmt.Println("\nPurge complete.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(purgeCmd)
}
