package kube

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeleteResource(client *kubernetes.Clientset, namespace string, res Resource) error {
	ctx := context.Background()
	switch res.Kind {
	case "Deployment":
		return client.AppsV1().Deployments(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "StatefulSet":
		return client.AppsV1().StatefulSets(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "DaemonSet":
		return client.AppsV1().DaemonSets(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "Service":
		return client.CoreV1().Services(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "Job":
		return client.BatchV1().Jobs(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "CronJob":
		return client.BatchV1().CronJobs(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "PVC":
		return client.CoreV1().PersistentVolumeClaims(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "Secret":
		return client.CoreV1().Secrets(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	case "ConfigMap":
		return client.CoreV1().ConfigMaps(namespace).Delete(ctx, res.Name, metav1.DeleteOptions{})
	default:
		return fmt.Errorf("unsupported kind: %s", res.Kind)
	}
}
