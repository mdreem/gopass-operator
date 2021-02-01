package gopass_repository

import (
	"context"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"log"
)

type RepositoryServer struct {
}

func (*RepositoryServer) InitializeRepository(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	log.Printf("InitializeRepository called with: %v", *repository)
	return nil, nil
}

func (*RepositoryServer) UpdateRepository(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	log.Printf("UpdateRepository called with: %v", *repository)
	return nil, nil
}

func Initialize() *RepositoryServer {
	return &RepositoryServer{}
}
