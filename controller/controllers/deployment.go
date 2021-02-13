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
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var getRelevantDeploymentFunc = getRelevantDeployment

func (r *GopassRepositoryReconciler) createRepositoryServer(ctx context.Context, log logr.Logger, namespacedName types.NamespacedName) (bool, error) {
	var deployment *appsv1.Deployment
	deployment, err := r.getDeployment(ctx, log, namespacedName)
	if err != nil {
		log.Error(err, "unable to fetch deployment")
		return false, err
	}

	if deployment == nil {
		log.Info("creating deployment")

		deployment = r.createDeployment(namespacedName)
		err := r.Client.Create(ctx, deployment)
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

func (r *GopassRepositoryReconciler) getDeployment(ctx context.Context, log logr.Logger, namespacedName types.NamespacedName) (*appsv1.Deployment, error) {
	labelSelector := labels.Set{
		"gopassRepoName":      namespacedName.Name,
		"gopassRepoNamespace": namespacedName.Namespace,
	}

	var deployments = &appsv1.DeploymentList{}
	err := r.Client.List(ctx, deployments, &client.ListOptions{
		LabelSelector: labelSelector.AsSelector(),
		Namespace:     r.Namespace,
	})
	if err != nil {
		log.Error(err, "unable to fetch list of deployments")
		return nil, err
	}

	return getRelevantDeploymentFunc(&deployments.Items)
}

func getRelevantDeployment(deployments *[]appsv1.Deployment) (*appsv1.Deployment, error) {
	if len(*deployments) == 0 {
		return nil, nil
	}

	if len(*deployments) != 1 {
		return nil, fmt.Errorf("expected 1 deployment, found: %d", len(*deployments))
	}

	return &(*deployments)[0], nil
}

func (r *GopassRepositoryReconciler) createDeployment(namespacedName types.NamespacedName) *appsv1.Deployment {
	appName := namespacedName.Name + "-" + uuid.New().String()

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    r.Namespace,
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

func (r *GopassRepositoryReconciler) deleteDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	err := r.Client.Delete(ctx, deployment)
	return err
}
