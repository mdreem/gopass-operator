package gopass_repository

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/mdreem/gopass-operator/gopass-server/gopass_repository/cluster"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"k8s.io/apimachinery/pkg/types"
	"log"
)

type RepositoryServer struct {
	Repositories map[string]*gopassRepo
	Client       cluster.Client
}

type gopassRepo struct {
	store      *api.Gopass
	directory  string
	repository *git.Repository
}

func Initialize() (*RepositoryServer, error) {
	someClient, err := cluster.New()
	if err != nil {
		log.Printf("unable to create client: %v", err)
		return nil, err
	}

	return &RepositoryServer{
		Repositories: make(map[string]*gopassRepo),
		Client:       &someClient,
	}, nil
}

func (r *RepositoryServer) InitializeRepository(ctx context.Context, repositoryInitialization *gopass_repository.RepositoryInitialization) (*gopass_repository.RepositoryResponse, error) {
	err := r.initializeRepository(ctx, repositoryInitialization)
	if err != nil {
		return &gopass_repository.RepositoryResponse{
			Successful:   false,
			ErrorMessage: fmt.Sprintf("unable to initialize repository: %v", err),
		}, err
	}

	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}
func (r *RepositoryServer) UpdateRepository(ctx context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	err := r.updateRepository(ctx, repository)

	if err != nil {
		return &gopass_repository.RepositoryResponse{
			Successful:   false,
			ErrorMessage: fmt.Sprintf("unable to update repository: %v", err),
		}, err
	}

	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}
func (r *RepositoryServer) UpdateAllPasswords(ctx context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	err := r.updateAllPasswords(ctx, repository)

	if err != nil {
		return &gopass_repository.RepositoryResponse{
			Successful:   false,
			ErrorMessage: fmt.Sprintf("unable to update passwords: %v", err),
		}, err
	}

	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *RepositoryServer) DeleteSecret(ctx context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	successful, err := deleteSecretMap(ctx, types.NamespacedName{
		Namespace: repository.SecretName.Namespace,
		Name:      repository.SecretName.Name,
	})

	return &gopass_repository.RepositoryResponse{
		Successful:   successful,
		ErrorMessage: fmt.Sprintf("failed to delete Secret: %v", err),
	}, nil
}
