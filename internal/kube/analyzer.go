package kube

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubeClient() (*kubernetes.Clientset, error) {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func AnalyzeNamespace(client *kubernetes.Clientset, namespace string, rules AnalyzerRules) ([]Resource, error) {
	var results []Resource
	ctx := context.Background()
	now := time.Now()

	referencedSecrets := map[string]bool{}
	referencedConfigMaps := map[string]bool{}

	// Deployments
	deployments, _ := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	for _, d := range deployments.Items {
		addResourceRefs(&d.Spec.Template.Spec, referencedSecrets, referencedConfigMaps)
		age := int(now.Sub(d.CreationTimestamp.Time).Hours() / 24)
		isStale := false
		if rules.DeploymentReplicas != nil && *rules.DeploymentReplicas >= 0 && d.Status.Replicas <= int32(*rules.DeploymentReplicas) {
			isStale = true
		}
		if rules.OlderThanDays != nil && *rules.OlderThanDays >= 0 && age >= *rules.OlderThanDays {
			isStale = true
		}
		if isStale {
			results = append(results, Resource{"Deployment", d.Name})
		}
	}

	// StatefulSets
	ssets, _ := client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	for _, s := range ssets.Items {
		addResourceRefs(&s.Spec.Template.Spec, referencedSecrets, referencedConfigMaps)
		age := int(now.Sub(s.CreationTimestamp.Time).Hours() / 24)
		isStale := false
		if rules.OlderThanDays != nil && *rules.OlderThanDays >= 0 && age >= *rules.OlderThanDays {
			isStale = true
		}
		if isStale {
			results = append(results, Resource{"StatefulSet", s.Name})
		}
	}

	// DaemonSets
	daemonsets, _ := client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	for _, ds := range daemonsets.Items {
		addResourceRefs(&ds.Spec.Template.Spec, referencedSecrets, referencedConfigMaps)
		age := int(now.Sub(ds.CreationTimestamp.Time).Hours() / 24)
		if rules.OlderThanDays != nil && *rules.OlderThanDays >= 0 && age >= *rules.OlderThanDays {
			results = append(results, Resource{"DaemonSet", ds.Name})
		}
	}

	// Services
	services, _ := client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	for _, svc := range services.Items {
		isOrphan := false

		// Only consider services with selectors (ignore headless or ExternalName)
		if len(svc.Spec.Selector) > 0 {
			// Build label selector string
			sel := metav1.LabelSelector{MatchLabels: svc.Spec.Selector}
			selector, err := metav1.LabelSelectorAsSelector(&sel)
			if err != nil {
				continue
			}
			pods, _ := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{LabelSelector: selector.String()})
			if len(pods.Items) == 0 {
				isOrphan = true
			}
		} else {
			// No selector means it can't match pods; treat as orphan if desired
			isOrphan = true
		}

		age := int(now.Sub(svc.CreationTimestamp.Time).Hours() / 24)
		ageOk := rules.OlderThanDays == nil || *rules.OlderThanDays < 0 || age >= *rules.OlderThanDays
		orphanOk := rules.OrphanedServicesOnly == nil || *rules.OrphanedServicesOnly == false || isOrphan

		if orphanOk && ageOk {
			results = append(results, Resource{"Service", svc.Name})
		}
	}

	// CronJobs
	cronjobs, _ := client.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})
	for _, cj := range cronjobs.Items {
		age := int(now.Sub(cj.CreationTimestamp.Time).Hours() / 24)
		if rules.OlderThanDays != nil && *rules.OlderThanDays >= 0 && age >= *rules.OlderThanDays {
			results = append(results, Resource{"CronJob", cj.Name})
		}
	}

	// Jobs
	jobs, _ := client.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})
	for _, job := range jobs.Items {
		age := int(now.Sub(job.CreationTimestamp.Time).Hours() / 24)
		finished := job.Status.Succeeded > 0 || job.Status.Failed > 0
		if rules.OlderThanDays != nil && *rules.OlderThanDays >= 0 && age < *rules.OlderThanDays {
			continue
		}
		if rules.JobCompletedOnly != nil && *rules.JobCompletedOnly && !finished {
			continue
		}
		results = append(results, Resource{"Job", job.Name})
	}

	// PVCs
	pvcs, _ := client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	for _, pvc := range pvcs.Items {
		phase := string(pvc.Status.Phase)
		if len(rules.PVCPhases) > 0 {
			for _, phaseMatch := range rules.PVCPhases {
				if strings.EqualFold(phase, phaseMatch) {
					results = append(results, Resource{"PVC", pvc.Name})
					break
				}
			}
		} else if rules.OlderThanDays != nil && *rules.OlderThanDays >= 0 {
			age := int(now.Sub(pvc.CreationTimestamp.Time).Hours() / 24)
			if age >= *rules.OlderThanDays {
				results = append(results, Resource{"PVC", pvc.Name})
			}
		}
	}

	// Secrets
	secrets, _ := client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})
	for _, s := range secrets.Items {
		typeAllowed := true
		if len(rules.SecretTypes) > 0 {
			typeAllowed = containsString(rules.SecretTypes, string(s.Type))
		}
		orphaned := !referencedSecrets[s.Name]
		orphanOk := rules.OrphanedSecretsOnly == nil || *rules.OrphanedSecretsOnly == false || orphaned
		age := int(now.Sub(s.CreationTimestamp.Time).Hours() / 24)
		ageOk := rules.OlderThanDays == nil || *rules.OlderThanDays < 0 || age >= *rules.OlderThanDays
		if typeAllowed && orphanOk && ageOk {
			results = append(results, Resource{"Secret", s.Name})
		}
	}

	// ConfigMaps
	configmaps, _ := client.CoreV1().ConfigMaps(namespace).List(ctx, metav1.ListOptions{})
	for _, c := range configmaps.Items {
		orphaned := !referencedConfigMaps[c.Name]
		orphanOk := rules.OrphanedConfigMapsOnly == nil || *rules.OrphanedConfigMapsOnly == false || orphaned
		age := int(now.Sub(c.CreationTimestamp.Time).Hours() / 24)
		ageOk := rules.OlderThanDays == nil || *rules.OlderThanDays < 0 || age >= *rules.OlderThanDays
		if orphanOk && ageOk {
			results = append(results, Resource{"ConfigMap", c.Name})
		}
	}

	return results, nil
}

// Helper: Scan PodSpec for references to Secrets/ConfigMaps
func addResourceRefs(spec *corev1.PodSpec, secrets, configmaps map[string]bool) {
	for _, vol := range spec.Volumes {
		if vol.Secret != nil {
			secrets[vol.Secret.SecretName] = true
		}
		if vol.ConfigMap != nil {
			configmaps[vol.ConfigMap.Name] = true
		}
	}
	for _, c := range spec.Containers {
		for _, ef := range c.EnvFrom {
			if ef.ConfigMapRef != nil && ef.ConfigMapRef.Name != "" {
				configmaps[ef.ConfigMapRef.Name] = true
			}
			if ef.SecretRef != nil && ef.SecretRef.Name != "" {
				secrets[ef.SecretRef.Name] = true
			}
		}
		for _, e := range c.Env {
			if e.ValueFrom != nil {
				if e.ValueFrom.ConfigMapKeyRef != nil {
					configmaps[e.ValueFrom.ConfigMapKeyRef.Name] = true
				}
				if e.ValueFrom.SecretKeyRef != nil {
					secrets[e.ValueFrom.SecretKeyRef.Name] = true
				}
			}
		}
	}
}

func containsString(list []string, v string) bool {
	for _, s := range list {
		if strings.EqualFold(s, v) {
			return true
		}
	}
	return false
}
