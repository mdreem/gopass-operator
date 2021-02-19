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
	"google.golang.org/grpc"
	"time"

	"github.com/go-logr/logr"
	gopassv1alpha1 "github.com/mdreem/gopass-operator/controller/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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

	service, err := r.getService(ctx, req.NamespacedName)

	// due to the creation/deletion-logic there are cases where the Service does not exist.
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

	result, err, done := r.handleDeletionOfResource(ctx, req, gopassRepository, repositoryServiceClient)
	if done {
		return result, err
	}

	deploymentFinished, err := r.createRepositoryServer(ctx, req.NamespacedName)
	if err != nil {
		log.Error(err, "not able to deploy repository server")
		return ctrl.Result{}, err
	}

	if !deploymentFinished {
		log.Info("deployment not yet ready, trying again later")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	if service == nil {
		log.Info("service did not exist yet, trying again later")
		return ctrl.Result{RequeueAfter: 5 * time.Second}, nil
	}

	err = initializeRepository(ctx, log, gopassRepository.Spec.RepositoryURL, repositoryServiceClient, req.Namespace, gopassRepository.Spec)
	if err != nil {
		log.Error(err, "unable to initialize repository")
		return ctrl.Result{}, err
	}

	err = updateRepository(ctx, req, repositoryServiceClient, gopassRepository)
	if err != nil {
		log.Error(err, "unable to update repository")
		return ctrl.Result{}, err
	}

	err = updateAllPasswords(ctx, log, req.NamespacedName, gopassRepository.Spec.RepositoryURL, repositoryServiceClient)
	if err != nil {
		log.Error(err, "unable to fetch secrets")
		return ctrl.Result{}, err
	}

	interval, err := parseRefreshInterval(gopassRepository.Spec.RefreshInterval)
	if err != nil {
		log.Error(err, "unable to parse refresh interval")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: interval}, nil
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

// SetupWithManager sets up the controller with the Manager.
func (r *GopassRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gopassv1alpha1.GopassRepository{}).
		Complete(r)
}
