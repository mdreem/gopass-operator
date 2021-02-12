package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func createRepositoryServer(ctx context.Context, log logr.Logger, namespacedName types.NamespacedName) (bool, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return false, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}

	deployment, err := getDeployment(ctx, clientset, namespacedName)
	if err != nil {
		log.Error(err, "unable to fetch deployment")
		return false, err
	}

	if deployment == nil {
		log.Info("creating deployment")

		newDeployment := createDeployment(namespacedName)
		_, err := clientset.AppsV1().Deployments("operator-system").Create(ctx, newDeployment, metav1.CreateOptions{})
		if err != nil {
			log.Error(err, "unable to create deployment")
			return false, err
		}
	}

	availableReplicas := deployment.Status.AvailableReplicas
	if availableReplicas > 0 {
		log.Info("deployment of repository server finished")
		return true, nil
	}

	log.Info("deployment of repository server in progress")

	return false, nil
}

func getDeployment(ctx context.Context, clientset *kubernetes.Clientset, namespacedName types.NamespacedName) (*appsv1.Deployment, error) {
	labelSelector := labels.Set{
		"gopassRepoName":      namespacedName.Name,
		"gopassRepoNamespace": namespacedName.Namespace,
	}

	deployments, err := clientset.AppsV1().Deployments("operator-system").List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector.AsSelector().String(),
	})
	if err != nil {
		return nil, err
	}

	if len(deployments.Items) == 0 {
		return nil, nil
	}

	if len(deployments.Items) != 1 {
		return nil, fmt.Errorf("expected 1 deployment, found: %d", len(deployments.Items))
	}

	return &deployments.Items[0], nil
}

func createDeployment(namespacedName types.NamespacedName) *appsv1.Deployment {
	appName := namespacedName.Name + "-" + uuid.New().String()

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    "operator-system",
			GenerateName: namespacedName.Name + "-",
			Labels: map[string]string{
				"app":                 appName,
				"gopassRepoName":      namespacedName.Name,
				"gopassRepoNamespace": namespacedName.Namespace,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": appName},
			},
			Replicas: getIntPointer(1),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": appName},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  namespacedName.Name,
							Image: "gopass-server:latest",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9000,
								},
							},
							ImagePullPolicy: "Never",
						},
					},
				},
			},
		},
	}

	return deployment
}

func getIntPointer(val int32) *int32 {
	return &val
}
