package cluster

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Secret struct {
	Name     string
	Password string
}

type Client interface {
	GetRepositoryCredentials(ctx context.Context, authentication *gopass_repository.Authentication) (Secret, error)
	GetGpgKey(ctx context.Context, authentication *gopass_repository.Authentication) error
}

type KubernetesClient struct {
	clientset *kubernetes.Clientset
}

func New() (KubernetesClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return KubernetesClient{}, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return KubernetesClient{}, err
	}

	return KubernetesClient{
		clientset: clientset,
	}, nil
}

func (k *KubernetesClient) GetRepositoryCredentials(ctx context.Context, authentication *gopass_repository.Authentication) (Secret, error) {
	secretMap, err := k.clientset.CoreV1().Secrets((*authentication).Namespace).Get(ctx, authentication.SecretRef, metav1.GetOptions{})
	if err != nil {
		return Secret{}, err
	}

	password, ok := (*secretMap).Data[authentication.SecretKey]
	if !ok {
		return Secret{}, fmt.Errorf("unable to find key '%s' in secret '%s' in namespace '%s'", authentication.SecretKey, authentication.SecretRef, authentication.Namespace)
	}

	return Secret{
		Name:     authentication.Username,
		Password: string(password),
	}, nil
}

func (k *KubernetesClient) GetGpgKey(ctx context.Context, authentication *gopass_repository.Authentication) error {
	log.Printf("add gpg key")

	secretMap, err := k.clientset.CoreV1().Secrets((*authentication).Namespace).Get(ctx, "gpg-key", metav1.GetOptions{})
	if err != nil {
		return err
	}

	gpgKey, ok := (*secretMap).Data["gpg-key"]
	if !ok {
		return fmt.Errorf("unable to find key '%s' in secret '%s' in namespace '%s'", authentication.SecretKey, authentication.SecretRef, authentication.Namespace)
	}

	_, err = addKey(ctx, gpgKey)
	if err != nil {
		return err
	}

	return nil
}

func addKey(ctx context.Context, key []byte) ([]byte, error) {
	args := make([]string, 0)
	args = append(args, "--import")
	cmd := exec.CommandContext(ctx, "gpg", args...)
	cmd.Stdin = bytes.NewReader(key)
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	return cmd.Output()
}
