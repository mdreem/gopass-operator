/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"
	"time"

	"github.com/go-logr/logr"
	gopassv1alpha1 "github.com/mdreem/gopass-operator/controller/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const finalizerName = "gopass.repository.finalizer"

var createRepositoryServiceClientFunc = createRepositoryServiceClient

// GopassRepositoryReconciler reconciles a GopassRepository object
type GopassRepositoryReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	Namespace string
}

// +kubebuilder:rbac:groups=gopass.gopass.operator,resources=gopassrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gopass.gopass.operator,resources=gopassrepositories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gopass.gopass.operator,resources=gopassrepositories/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GopassRepository object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *GopassRepositoryReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("gopassrepository", req.NamespacedName)

	log.Info("called reconcile for gopassRepository")

	gopassRepository := &gopassv1alpha1.GopassRepository{}
	err := r.Get(ctx, req.NamespacedName, gopassRepository)
	if err != nil {
		log.Error(err, "unable to fetch GopassRepository from request")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	service, err := r.getService(ctx, log, req.NamespacedName)

	var repositoryServiceClient gopass_repository.RepositoryServiceClient
	if service != nil {
		var conn *grpc.ClientConn
		repositoryServiceClient, conn, err = createRepositoryServiceClientFunc(service.Name)
		if err != nil {
			log.Error(err, "not able to connect to repository server")
			return ctrl.Result{}, err
		}
		defer closeConnection(log, conn)
	}

	result, err, done := r.handleDeletion(ctx, req, log, gopassRepository, repositoryServiceClient)
	if done {
		return result, err
	}

	deploymentFinished, err := r.createRepositoryServer(ctx, log, req.NamespacedName)
	if err != nil {
		log.Error(err, "not able to deploy repository server")
		return ctrl.Result{}, err
	}

	if !deploymentFinished {
		log.Info("deployment not yet ready")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	if service == nil {
		log.Info("service did not exist yet")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	err = initializeRepository(ctx, log, gopassRepository.Spec.RepositoryURL, repositoryServiceClient, req.Namespace, gopassRepository.Spec)
	if err != nil {
		log.Error(err, "unable to initialize repository")
		return ctrl.Result{}, err
	}

	interval, err := parseRefreshInterval(gopassRepository.Spec.RefreshInterval)
	if err != nil {
		log.Error(err, "unable to parse refresh interval")
		return ctrl.Result{}, err
	}

	_, err = repositoryServiceClient.UpdateRepository(ctx, &gopass_repository.Repository{
		RepositoryURL: gopassRepository.Spec.RepositoryURL,
		Authentication: &gopass_repository.Authentication{
			Namespace: req.NamespacedName.Namespace,
			Username:  gopassRepository.Spec.UserName,
			SecretRef: gopassRepository.Spec.SecretKeyRef.Name,
			SecretKey: gopassRepository.Spec.SecretKeyRef.Key,
		},
	})
	if err != nil {
		log.Error(err, "unable to update repository")
		return ctrl.Result{}, err
	}

	err = updateAllPasswords(ctx, log, req.NamespacedName, gopassRepository.Spec.RepositoryURL, repositoryServiceClient)
	if err != nil {
		log.Error(err, "unable to fetch secrets")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: interval}, nil
}

func (r *GopassRepositoryReconciler) handleDeletion(ctx context.Context, req ctrl.Request, log logr.Logger, repository *gopassv1alpha1.GopassRepository, serviceClient gopass_repository.RepositoryServiceClient) (ctrl.Result, error, bool) {
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

func closeConnection(log logr.Logger, conn *grpc.ClientConn) {
	connectionError := conn.Close()
	if connectionError != nil {
		log.Error(connectionError, "not able to close connection")
	}
	log.Info("closed connection to repository server")
}

func parseRefreshInterval(refreshInterval string) (time.Duration, error) {
	if refreshInterval == "" {
		return 0, nil
	}

	refreshIntervalValue, err := time.ParseDuration(refreshInterval)
	if err != nil {
		return 0, err
	}
	return refreshIntervalValue, nil
}

func createRepositoryServiceClient(targetUrl string) (gopass_repository.RepositoryServiceClient, *grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(targetUrl+":9000", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return gopass_repository.NewRepositoryServiceClient(conn), conn, nil
}

func updateAllPasswords(ctx context.Context, log logr.Logger, namespacedName types.NamespacedName, url string, repositoryServiceClient gopass_repository.RepositoryServiceClient) error {
	_, err := repositoryServiceClient.UpdateAllPasswords(ctx,
		&gopass_repository.Repository{
			RepositoryURL: url,
			SecretName: &gopass_repository.NamespacedName{
				Namespace: namespacedName.Namespace,
				Name:      namespacedName.Name,
			},
		})

	if err != nil {
		log.Error(err, "not able to fetch passwords")
		return err
	}
	return nil
}

func initializeRepository(ctx context.Context, log logr.Logger, url string, repositoryServiceClient gopass_repository.RepositoryServiceClient,
	namespace string, gopassRepositorySpec gopassv1alpha1.GopassRepositorySpec) error {
	log.Info("attempting to call repository server")
	repository, err := repositoryServiceClient.InitializeRepository(
		ctx,
		&gopass_repository.RepositoryInitialization{
			Repository: &gopass_repository.Repository{
				RepositoryURL: url,
				Authentication: &gopass_repository.Authentication{
					Namespace: namespace,
					Username:  gopassRepositorySpec.UserName,
					SecretRef: gopassRepositorySpec.SecretKeyRef.Name,
					SecretKey: gopassRepositorySpec.SecretKeyRef.Key,
				},
			},
			GpgKeyReference: &gopass_repository.GpgKeyReference{
				GpgKeyRef:    gopassRepositorySpec.GpgKeyRef.Name,
				GpgKeyRefKey: gopassRepositorySpec.GpgKeyRef.Key,
			},
		},
	)

	if err != nil {
		log.Error(err, "invalid response")
		return err
	}

	if repository != nil {
		log.Info("repository call:", "successful", (*repository).Successful)
	} else {
		log.Info("empty response from repository server")
	}

	return nil
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

// SetupWithManager sets up the controller with the Manager.
func (r *GopassRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gopassv1alpha1.GopassRepository{}).
		Complete(r)
}
