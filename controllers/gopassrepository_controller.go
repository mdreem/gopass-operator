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
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"regexp"
	"time"

	"github.com/go-logr/logr"
	gopassv1alpha1 "github.com/mdreem/gopass-operator/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	corev1 "k8s.io/api/core/v1"
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

	repositoryServiceClient, conn, err := createRepositoryServiceClient()
	if err != nil {
		log.Error(err, "not able to connect to repository server")
		return ctrl.Result{}, err
	}
	defer closeConnection(log, conn)

	gopassRepository := &gopassv1alpha1.GopassRepository{}
	err = r.Get(ctx, req.NamespacedName, gopassRepository)
	if err != nil {
		log.Error(err, "unable to fetch data from request")
		return ctrl.Result{}, err
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

	secrets, err := fetchAllPasswords(ctx, log, gopassRepository.Spec.RepositoryURL, repositoryServiceClient)
	if err != nil {
		log.Error(err, "unable to fetch secrets")
		return ctrl.Result{}, err
	}

	errUpdateSecretMap := r.updateSecretMap(ctx, log, req.NamespacedName, secrets)
	if errUpdateSecretMap != nil {
		log.Error(errUpdateSecretMap, "error updating secret map")
		return ctrl.Result{}, errUpdateSecretMap
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

func createRepositoryServiceClient() (gopass_repository.RepositoryServiceClient, *grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("operator-gopass-repository:9000", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return gopass_repository.NewRepositoryServiceClient(conn), conn, nil
}

func fetchAllPasswords(ctx context.Context, log logr.Logger, url string, repositoryServiceClient gopass_repository.RepositoryServiceClient) (*gopass_repository.SecretList, error) {
	passwords, err := repositoryServiceClient.FetchAllPasswords(ctx,
		&gopass_repository.Repository{
			RepositoryURL: url,
		})

	if err != nil {
		log.Error(err, "not able to fetch passwords")
		return nil, err
	}
	return passwords, nil
}

func initializeRepository(ctx context.Context, log logr.Logger, url string, repositoryServiceClient gopass_repository.RepositoryServiceClient,
	namespace string, gopassRepositorySpec gopassv1alpha1.GopassRepositorySpec) error {
	log.Info("attempting to call repository server")
	repository, responseErr := repositoryServiceClient.InitializeRepository(
		ctx,
		&gopass_repository.Repository{
			RepositoryURL: url,
			Authentication: &gopass_repository.Authentication{
				Namespace: namespace,
				Username:  gopassRepositorySpec.UserName,
				SecretRef: gopassRepositorySpec.SecretKeyRef.Name,
				SecretKey: gopassRepositorySpec.SecretKeyRef.Key,
			},
		},
	)

	passwords, err := repositoryServiceClient.FetchAllPasswords(ctx, &gopass_repository.Repository{
		RepositoryURL: url,
	})
	if err != nil {
		log.Error(err, "not able to fetch passwords")
		return err
	}

	for idx, password := range passwords.Secrets {
		log.Info("password", "index", idx, "name", password.Name, "password", password.Password)
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

func (r *GopassRepositoryReconciler) updateSecretMap(ctx context.Context, log logr.Logger, namespacedName types.NamespacedName, secrets *gopass_repository.SecretList) error {
	log.Info("updating secret map")

	newSecret := createSecret(secrets, namespacedName)

	secret := corev1.Secret{}
	objectKey := client.ObjectKey{
		Name:      namespacedName.Name,
		Namespace: namespacedName.Namespace,
	}

	err := r.Client.Get(ctx, objectKey, &secret)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("creating secret map")

			err := r.Create(ctx, &newSecret)
			if err != nil {
				log.Error(err, "not able to create secret map")
				return err
			}
			return nil
		} else {
			log.Error(err, "not able to fetch secret map")
			return err
		}
	}

	secret.StringData = createSecretMap(secrets)
	err = r.Update(ctx, &secret)
	if err != nil {
		log.Error(err, "not able to update secret map")
		return err
	}

	return nil
}

func createSecret(secrets *gopass_repository.SecretList, namespacedName types.NamespacedName) corev1.Secret {
	newSecretMap := createSecretMap(secrets)

	newSecret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      namespacedName.Name,
			Namespace: namespacedName.Namespace,
		},
		StringData: newSecretMap,
	}
	return newSecret
}

func createSecretMap(secrets *gopass_repository.SecretList) map[string]string {
	newSecretMap := make(map[string]string)
	for _, secret := range secrets.Secrets {
		newSecretMap[rename(secret.Name)] = secret.Password
	}
	return newSecretMap
}

func rename(name string) string {
	reg, _ := regexp.Compile("[^a-zA-Z0-9 ]+")
	return reg.ReplaceAllString(name, "-")
}
