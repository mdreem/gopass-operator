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
	gopassv1alpha1 "github.com/mdreem/gopass-operator/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// GopassRepositoryReconciler reconciles a GopassRepository object
type GopassRepositoryReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gopass.gopass.operator,resources=gopassrepositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gopass.gopass.operator,resources=gopassrepositories/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gopass.gopass.operator,resources=gopassrepositories/finalizers,verbs=update

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

	log.Info("reconciliation")

	gopassRepository := &gopassv1alpha1.GopassRepository{}
	err := r.Get(ctx, req.NamespacedName, gopassRepository)
	if err != nil {
		log.Error(err, "unable to fetch data from request")
		return ctrl.Result{}, err
	}

	err = r.initializeRepository(ctx, log, gopassRepository.Spec.RepositoryURL)
	if err != nil {
		log.Error(err, "unable to initialize repository")
		return ctrl.Result{}, err
	}

	interval, err := parseRefreshInterval(gopassRepository.Spec.RefreshInterval)
	if err != nil {
		log.Error(err, "unable to parse refresh interval")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: interval}, nil
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

func (r *GopassRepositoryReconciler) initializeRepository(ctx context.Context, log logr.Logger, url string) error {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("operator-gopass-repository:9000", grpc.WithInsecure())
	if err != nil {
		log.Error(err, "not able to connect to repository server")
		return err
	}

	log.Info("attempting to call repository server")

	c := gopass_repository.NewRepositoryServiceClient(conn)

	repository, responseErr := c.InitializeRepository(
		ctx,
		&gopass_repository.Repository{
			RepositoryURL: url,
		},
	)

	passwords, passwordFetchErr := c.FetchAllPasswords(ctx, &gopass_repository.Repository{
		RepositoryURL: url,
	})
	if passwordFetchErr != nil {
		log.Error(err, "not able to fetch passwords")
		return err
	}

	for idx, password := range passwords.Secrets {
		log.Info("password", "index", idx, "name", password.Name, "password", password.Password)
	}

	connectionError := conn.Close()
	if connectionError != nil {
		log.Error(connectionError, "not able to close connection")
		return connectionError
	}

	if responseErr != nil {
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

// SetupWithManager sets up the controller with the Manager.
func (r *GopassRepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gopassv1alpha1.GopassRepository{}).
		Complete(r)
}
