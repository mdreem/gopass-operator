// +build mock

package cluster

import (
	"context"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
)

type KubernetesTestClient struct {
}

func (*KubernetesTestClient) GetRepositoryCredentials(_ context.Context, _ *gopass_repository.Authentication) (Secret, error) {
	return Secret{}, nil
}

func (*KubernetesTestClient) GetGpgKey(_ context.Context, _ *gopass_repository.Authentication) error {
	return nil
}
