package gopass_repository

import (
	"context"
	"fmt"
	"github.com/mdreem/gopass-operator/gopass-server/gopass_repository/cluster"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"regexp"
)

func (r *RepositoryServer) updateAllPasswords(ctx context.Context, repository *gopass_repository.Repository) error {
	repo, ok := (r.Repositories)[(*repository).RepositoryURL]
	if !ok {
		return fmt.Errorf("repository with URL '%s' not found", (*repository).RepositoryURL)
	}

	passwords, err := fetchAllPasswords(ctx, repo)
	if err != nil {
		log.Printf("error fetching passwords: %v\n", err)
		return err
	}

	secretList := gopass_repository.SecretList{
		Secrets: make([]*gopass_repository.Secret, 0),
	}

	for _, password := range passwords {
		secretList.Secrets = append(secretList.Secrets, &gopass_repository.Secret{
			Name:     password.Name,
			Password: password.Password,
		})
	}

	err = updateSecretMap(ctx, types.NamespacedName{
		Namespace: repository.SecretName.Namespace,
		Name:      repository.SecretName.Name,
	}, &secretList)
	if err != nil {
		log.Printf("unable to update secret map: %v\n", err)
		return err
	}

	return nil
}

func fetchAllPasswords(ctx context.Context, repo *gopassRepo) ([]cluster.Secret, error) {
	list, err := (*repo).store.List(ctx)
	if err != nil {
		log.Printf("not able to list contents of repository: %v\n", err)
		return nil, err
	}

	passwords := make([]cluster.Secret, 0)

	for _, passwordName := range list {
		password, err := (*repo).store.Get(ctx, passwordName, "")
		if err != nil {
			log.Printf("not able to fetch password '%s': %v\n", passwordName, err)
			continue
		}
		passwords = append(passwords, cluster.Secret{
			Name:     passwordName,
			Password: password.Password(),
		})
	}

	return passwords, nil
}

func updateSecretMap(ctx context.Context, namespacedName types.NamespacedName, secrets *gopass_repository.SecretList) error {
	log.Printf("updating secret map\n")

	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	newSecret := createSecret(secrets, namespacedName)

	_, err = getSecretMap(ctx, clientset, namespacedName)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Printf("creating secret map")

			_, err := clientset.CoreV1().Secrets(namespacedName.Namespace).Create(ctx, &newSecret, metav1.CreateOptions{})
			if err != nil {
				log.Printf("unable to create secret map: %v\n", err)
				return err
			}
			return nil
		} else {
			log.Printf("unable to fetch secret map: %v\n", err)
			return err
		}
	}

	_, err = clientset.CoreV1().Secrets(namespacedName.Namespace).Update(ctx, &newSecret, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("not able to update secret map: %v", err)
		return err
	}

	return nil
}

func getSecretMap(ctx context.Context, clientset *kubernetes.Clientset, namespacedName types.NamespacedName) (*corev1.Secret, error) {
	secretMap, err := clientset.CoreV1().Secrets(namespacedName.Namespace).Get(ctx, namespacedName.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return secretMap, nil
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
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(name, "-")
}

func deleteSecretMap(ctx context.Context, namespacedName types.NamespacedName) (bool, error) {
	log.Printf("deleting secret")
	config, err := rest.InClusterConfig()
	if err != nil {
		return false, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return false, err
	}

	secret, err := getSecretMap(ctx, clientset, namespacedName)

	if secret == nil {
		log.Printf("secret not found")
		return true, nil
	}

	err = clientset.CoreV1().Secrets(namespacedName.Namespace).Delete(ctx, namespacedName.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("unable to delete secret: %v", err)
		return false, err
	}

	return true, nil
}
