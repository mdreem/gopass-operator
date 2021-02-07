package gopass_repository

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	ssh_2 "golang.org/x/crypto/ssh"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type gopassRepo struct {
	store     *api.Gopass
	directory string
}

type config struct {
	Path string `yaml:"path"`
}

type RepositoryServer struct {
	Repositories map[string]*gopassRepo
}

type secret struct {
	Name     string
	Password string
}

func (r *RepositoryServer) InitializeRepository(ctx context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	log.Printf("InitializeRepository called with: %s", (*repository).RepositoryURL)

	_, ok := (r.Repositories)[(*repository).RepositoryURL]
	if ok {
		log.Printf("repository with URL '%s' already initialized", (*repository).RepositoryURL)
		return &gopass_repository.RepositoryResponse{
			Successful:   true,
			ErrorMessage: "",
		}, nil
	}

	credentials, err := getRepositoryCredentials(ctx, repository.Authentication)
	if err != nil {
		log.Printf("error initializing repository: %v", err)
		return nil, err
	}

	gopassRepository, err := initializeNewGopassRepository((*repository).RepositoryURL, credentials)
	if err != nil {
		log.Printf("error initializing repository: %v", err)
		return nil, err
	}

	(r.Repositories)[(*repository).RepositoryURL] = gopassRepository

	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (*RepositoryServer) UpdateRepository(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	log.Printf("UpdateRepository called with: %s", (*repository).RepositoryURL)
	return &gopass_repository.RepositoryResponse{
		Successful:   true,
		ErrorMessage: "",
	}, nil
}

func (r *RepositoryServer) FetchAllPasswords(ctx context.Context, repository *gopass_repository.Repository) (*gopass_repository.SecretList, error) {
	repo, ok := (r.Repositories)[(*repository).RepositoryURL]
	if !ok {
		return nil, fmt.Errorf("repository with URL '%s' not found", (*repository).RepositoryURL)
	}
	passwords, err := fetchAllPasswords(ctx, repo)
	if err != nil {
		log.Printf("error fetching passwords: %v", err)
		return nil, err
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

	return &secretList, nil
}

func Initialize() *RepositoryServer {
	return &RepositoryServer{
		Repositories: make(map[string]*gopassRepo),
	}
}

func initializeNewGopassRepository(repositoryUrl string, credentials secret) (*gopassRepo, error) {
	repoDir, err := ioutil.TempDir("", "gopass")
	if err != nil {
		log.Printf("not able to create local repository directory: %v", err)
		return nil, err
	}

	err = cloneGopassRepo(repositoryUrl, repoDir, credentials.Name, credentials.Password)
	if err != nil {
		log.Printf("not able clone gopass repository with URL %s: %v", repositoryUrl, err)
		return nil, err
	}

	gr, err := createNewGopassClient(context.Background(), repoDir)
	if err != nil {
		log.Printf("not able to create new gopass client: %v", err)
		return nil, err
	}

	return gr, nil
}

func cloneGopassRepo(repositoryUrl string, path string, username string, password string) error {
	sshPassword := ssh.Password{
		User:     username,
		Password: password,
		HostKeyCallbackHelper: ssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh_2.InsecureIgnoreHostKey(),
		},
	}

	log.Printf("username: '%s' - password: '%s'\n", username, password)

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repositoryUrl,
		Depth:    1,
		Progress: os.Stdout,
		Auth:     &sshPassword,
	})
	return err
}

func createNewGopassClient(ctx context.Context, path string) (*gopassRepo, error) {
	file, err := ioutil.TempFile("", "*config.yml")
	if err != nil {
		log.Printf("not able to create temporary configuration file %v\n", err)
		return nil, err
	}
	defer removeFile(file)
	log.Printf("created temporary configuration file: %s\n", file.Name())

	c := config{
		Path: path,
	}

	marshalledConfig, err := yaml.Marshal(&c)
	if err != nil {
		log.Printf("not able to marshall configuration file: %v\n", err)
		return nil, err
	}

	_, err = file.Write(marshalledConfig)
	if err != nil {
		log.Printf("not able to write configuration file: %v\n", err)
		return nil, err
	}

	err = file.Close()
	if err != nil {
		log.Printf("not able to close configuration file: %v\n", err)
		return nil, err
	}

	err = os.Setenv("GOPASS_CONFIG", file.Name())
	if err != nil {
		log.Printf("not able to set environment variable: %v\n", err)
		return nil, err
	}

	store, err := api.New(ctx)
	if err != nil {
		log.Printf("not able to create a new gopass client: %v\n", err)
		return nil, err
	}

	gr := &gopassRepo{
		store:     store,
		directory: path,
	}

	return gr, nil
}

func removeFile(file *os.File) {
	err := os.Remove(file.Name())
	if err != nil {
		log.Printf("failed to remove file: %v\n", err)
	}
}

func fetchAllPasswords(ctx context.Context, repo *gopassRepo) ([]secret, error) {
	list, err := (*repo).store.List(ctx)
	if err != nil {
		log.Printf("not able to list contents of repository: %v\n", err)
		return nil, err
	}

	passwords := make([]secret, 0)

	for _, passwordName := range list {
		password, err := (*repo).store.Get(ctx, passwordName, "")
		if err != nil {
			log.Printf("not able to fetch password '%s': %v\n", passwordName, err)
			return nil, err
		}
		passwords = append(passwords, secret{
			Name:     passwordName,
			Password: password.Password(),
		})
	}

	return passwords, nil
}

func getRepositoryCredentials(ctx context.Context, authentication *gopass_repository.Authentication) (secret, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return secret{}, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return secret{}, err
	}

	secretMap, err := clientset.CoreV1().Secrets(authentication.Namespace).Get(ctx, authentication.SecretRef, metav1.GetOptions{})
	if err != nil {
		return secret{}, err
	}

	password, ok := (*secretMap).Data[authentication.SecretKey]
	if !ok {
		return secret{}, fmt.Errorf("unable to find key '%s' in secret '%s' in namespace '%s'", authentication.SecretKey, authentication.SecretRef, authentication.Namespace)
	}

	return secret{
		Name:     authentication.Username,
		Password: string(password),
	}, nil
}
