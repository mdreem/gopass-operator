package controllers

import (
	"context"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
)

type TestRepositoryServer struct {
	Calls map[string][]string
}

func InitializeTestRepositoryServer() *TestRepositoryServer {
	return &TestRepositoryServer{
		Calls: map[string][]string{
			"InitializeRepository": {},
			"UpdateRepository":     {},
			"UpdateAllPasswords":   {},
		},
	}
}

func (r *TestRepositoryServer) InitializeRepository(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	r.Calls["InitializeRepository"] = append(r.Calls["InitializeRepository"], repository.RepositoryURL)
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServer) UpdateRepository(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	r.Calls["UpdateRepository"] = append(r.Calls["UpdateRepository"], repository.RepositoryURL)
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServer) UpdateAllPasswords(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	r.Calls["UpdateAllPasswords"] = append(r.Calls["UpdateAllPasswords"], repository.RepositoryURL)
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *TestRepositoryServer) DeleteSecret(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	r.Calls["UpdateAllPasswords"] = append(r.Calls["UpdateAllPasswords"], repository.RepositoryURL)
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}
