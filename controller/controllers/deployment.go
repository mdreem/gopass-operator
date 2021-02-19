package controllers

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var getRelevantDeploymentFunc = getRelevantDeployment

func (r *GopassRepositoryReconciler) createRepositoryServer(ctx context.Context, namespacedName types.NamespacedName) (bool, error) {
	appName := namespacedName.Name + "-" + uuid.New().String()

	var deployment *appsv1.Deployment
	deployment, err := r.getDeployment(ctx, namespacedName)
	if err != nil {
		r.Log.Error(err, "unable to fetch deployment")
		return false, err
	}

	if deployment == nil {
		r.Log.Info("creating deployment")

		deployment = r.createDeployment(namespacedName, appName)
		err := r.Client.Create(ctx, deployment)
		if err != nil {
			r.Log.Error(err, "unable to create deployment")
			return false, err
		}
	}

	var service *corev1.Service
	service, err = r.getService(ctx, namespacedName)
	if err != nil {
		r.Log.Error(err, "unable to fetch service")
		return false, err
	}

	if service == nil {
		r.Log.Info("creating service")
		service = r.createService(namespacedName, appName)
		err := r.Client.Create(ctx, service)
		if err != nil {
			r.Log.Error(err, "unable to create service")
			return false, err
		}
	}

	availableReplicas := deployment.Status.AvailableReplicas
	if availableReplicas > 0 {
		r.Log.Info("deployment of repository server finished")
		return true, nil
	}

	r.Log.Info("deployment of repository server in progress")

	return false, nil
}

func (r *GopassRepositoryReconciler) getDeployment(ctx context.Context, namespacedName types.NamespacedName) (*appsv1.Deployment, error) {
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
		r.Log.Error(err, "unable to fetch list of deployments")
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

func (r *GopassRepositoryReconciler) getService(ctx context.Context, namespacedName types.NamespacedName) (*corev1.Service, error) {
	labelSelector := labels.Set{
		"gopassRepoName":      namespacedName.Name,
		"gopassRepoNamespace": namespacedName.Namespace,
	}

	var services = &corev1.ServiceList{}
	err := r.Client.List(ctx, services, &client.ListOptions{
		LabelSelector: labelSelector.AsSelector(),
		Namespace:     r.Namespace,
	})
	if err != nil {
		r.Log.Error(err, "unable to fetch list of services")
		return nil, err
	}

	return getRelevantService(&services.Items)
}

func getRelevantService(services *[]corev1.Service) (*corev1.Service, error) {
	if len(*services) == 0 {
		return nil, nil
	}

	if len(*services) != 1 {
		return nil, fmt.Errorf("expected 1 deployment, found: %d", len(*services))
	}

	return &(*services)[0], nil
}

func (r *GopassRepositoryReconciler) createDeployment(namespacedName types.NamespacedName, appName string) *appsv1.Deployment {
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

func (r *GopassRepositoryReconciler) createService(namespacedName types.NamespacedName, appName string) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace:    r.Namespace,
			GenerateName: namespacedName.Name + "-",
			Labels: map[string]string{
				"app":                 appName,
				"gopassRepoName":      namespacedName.Name,
				"gopassRepoNamespace": namespacedName.Namespace,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Protocol: "TCP",
					Port:     9000,
				},
			},
			Selector: map[string]string{"app": appName},
		},
	}

	return service
}

func getIntPointer(val int32) *int32 {
	return &val
}

func (r *GopassRepositoryReconciler) deleteDeployment(ctx context.Context, deployment *appsv1.Deployment) error {
	err := r.Client.Delete(ctx, deployment)
	return err
}

func (r *GopassRepositoryReconciler) deleteService(ctx context.Context, service *corev1.Service) error {
	err := r.Client.Delete(ctx, service)
	return err
}
