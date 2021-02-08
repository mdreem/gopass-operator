package gopass_repository

import (
	"context"
	"fmt"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
)

type RepositoryServer struct {
	Repositories map[string]*gopassRepo
}

type gopassRepo struct {
	store     *api.Gopass
	directory string
}

func (r *RepositoryServer) InitializeRepository(ctx context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	return r.initializeRepository(ctx, repository)
}
func (r *RepositoryServer) UpdateRepository(ctx context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	return r.updateRepository(ctx, repository)
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
