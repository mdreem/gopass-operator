package controllers

import (
	"context"
	"github.com/go-logr/logr"
	gopassv1alpha1 "github.com/mdreem/gopass-operator/controller/api/v1alpha1"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

func createRepositoryServiceClient(targetUrl string) (gopass_repository.RepositoryServiceClient, *grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(targetUrl+":9000", grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return gopass_repository.NewRepositoryServiceClient(conn), conn, nil
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

func updateRepository(ctx context.Context, req ctrl.Request, repositoryServiceClient gopass_repository.RepositoryServiceClient, gopassRepository *gopassv1alpha1.GopassRepository) error {
	_, err := repositoryServiceClient.UpdateRepository(ctx, &gopass_repository.Repository{
		RepositoryURL: gopassRepository.Spec.RepositoryURL,
		Authentication: &gopass_repository.Authentication{
			Namespace: req.NamespacedName.Namespace,
			Username:  gopassRepository.Spec.UserName,
			SecretRef: gopassRepository.Spec.SecretKeyRef.Name,
			SecretKey: gopassRepository.Spec.SecretKeyRef.Key,
		},
	})
	return err
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
