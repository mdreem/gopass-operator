package cluster

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var execCommandContext = exec.CommandContext

type Secret struct {
	Name     string
	Password string
}

type Client interface {
	GetRepositoryCredentials(ctx context.Context, authentication *gopass_repository.Authentication) (Secret, error)
	GetGpgKey(ctx context.Context, namespace string, gpgKeyReference *gopass_repository.GpgKeyReference) error
}

type KubernetesClient struct {
	clientset kubernetes.Interface
}

func New(kubernetesClient kubernetes.Interface) KubernetesClient {
	return KubernetesClient{
		clientset: kubernetesClient,
	}
}

func (k *KubernetesClient) GetRepositoryCredentials(ctx context.Context, authentication *gopass_repository.Authentication) (Secret, error) {
	secretMap, err := k.clientset.CoreV1().Secrets((*authentication).Namespace).Get(ctx, authentication.SecretRef, metav1.GetOptions{})
	if err != nil {
		log.Printf("unable to fetch Secret: %v", err)
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

func (k *KubernetesClient) GetGpgKey(ctx context.Context, namespace string, gpgKeyReference *gopass_repository.GpgKeyReference) error {
	log.Printf("add gpg key")

	secretMap, err := k.clientset.CoreV1().Secrets(namespace).Get(ctx, gpgKeyReference.GpgKeyRef, metav1.GetOptions{})
	if err != nil {
		log.Printf("unable to fetch Secret: %v", err)
		return err
	}

	gpgKey, ok := (*secretMap).Data[gpgKeyReference.GpgKeyRefKey]
	if !ok {
		return fmt.Errorf("unable to find key '%s' in secret '%s' in namespace '%s'", gpgKeyReference.GpgKeyRefKey, gpgKeyReference.GpgKeyRef, namespace)
	}

	_, err = addKey(ctx, gpgKey)
	if err != nil {
		log.Printf("unable to add key: %v", err)
		return err
	}

	return nil
}

func addKey(ctx context.Context, key []byte) ([]byte, error) {
	args := make([]string, 0)
	args = append(args, "--import")
	cmd := execCommandContext(ctx, "gpg", args...)
	cmd.Stdin = bytes.NewReader(key)
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	return cmd.Output()
}
