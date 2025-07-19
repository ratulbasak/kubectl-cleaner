package kube

type AnalyzerRules struct {
	DeploymentReplicas     *int     `yaml:"deployment_replicas"`
	JobCompletedOnly       *bool    `yaml:"job_completed_only"`
	PVCPhases              []string `yaml:"pvc_phases"`
	SecretTypes            []string `yaml:"secret_types"`
	OrphanedServicesOnly   *bool    `yaml:"orphaned_services_only"`
	OrphanedSecretsOnly    *bool    `yaml:"orphaned_secrets_only"`
	OrphanedConfigMapsOnly *bool    `yaml:"orphaned_configmaps_only"`
	OlderThanDays          *int     `yaml:"older_than"`
}

type Resource struct {
	Kind string
	Name string
}
