package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	gopassv1alpha1 "github.com/mdreem/gopass-operator/controller/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const finalizerName = "gopass.repository.finalizer"

func (r *GopassRepositoryReconciler) handleDeletionOfResource(ctx context.Context, req ctrl.Request, log logr.Logger, repository *gopassv1alpha1.GopassRepository, serviceClient gopass_repository.RepositoryServiceClient) (ctrl.Result, error, bool) {
	if repository.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("preparing to delete")
		if !containsString(repository.ObjectMeta.Finalizers, finalizerName) {
			repository.ObjectMeta.Finalizers = append(repository.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), repository); err != nil {
				return ctrl.Result{}, err, true
			}
		}
	} else {
		if containsString(repository.ObjectMeta.Finalizers, finalizerName) {
			err := r.deleteExternalResources(ctx, log, req.NamespacedName, serviceClient)
			if err != nil {
				return ctrl.Result{}, err, true
			}

			repository.ObjectMeta.Finalizers = removeString(repository.ObjectMeta.Finalizers, finalizerName)
			if err := r.Update(context.Background(), repository); err != nil {
				return ctrl.Result{}, err, true
			}
		}
		return ctrl.Result{}, nil, true
	}
	return ctrl.Result{}, nil, false
}

func (r *GopassRepositoryReconciler) deleteExternalResources(ctx context.Context, log logr.Logger, namespacedName types.NamespacedName, serviceClient gopass_repository.RepositoryServiceClient) error {
	if serviceClient != nil {
		secret, err := serviceClient.DeleteSecret(ctx, &gopass_repository.Repository{
			SecretName: &gopass_repository.NamespacedName{
				Namespace: namespacedName.Namespace,
				Name:      namespacedName.Name,
			},
		})
		if err != nil {
			log.Error(err, "unable to delete secret")
			return err
		}
		if !secret.Successful {
			delError := fmt.Errorf("deletion of secret not successful")
			log.Error(delError, "deleteExternalResourcesFailed")
			return delError
		}
	}

	deployment, err := r.getDeployment(ctx, log, namespacedName)
	if err != nil {
		log.Error(err, "unable to get deployment for deleteExternalResources")
		return err
	}

	if deployment != nil {
		err = r.deleteDeployment(ctx, deployment)
		if err != nil {
			log.Error(err, "unable to delete deployment")
			return err
		}
	}

	service, err := r.getService(ctx, log, namespacedName)
	if err != nil {
		log.Error(err, "unable to get deployment for deleteExternalResources")
		return err
	}

	if service != nil {
		err = r.deleteService(ctx, service)
		if err != nil {
			log.Error(err, "unable to delete service")
			return err
		}
	}

	return nil
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
