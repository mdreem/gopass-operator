package gopass_repository

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"io/ioutil"
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

func (r *RepositoryServer) InitializeRepository(_ context.Context, repository *gopass_repository.Repository) (*gopass_repository.RepositoryResponse, error) {
	log.Printf("InitializeRepository called with: %s", (*repository).RepositoryURL)

	_, ok := (r.Repositories)[(*repository).RepositoryURL]
	if ok {
		log.Printf("repository with URL '%s' already initialized", (*repository).RepositoryURL)
		return &gopass_repository.RepositoryResponse{
			Successful:   true,
			ErrorMessage: "",
		}, nil
	}

	gopassRepository, err := initializeNewGopassRepository((*repository).RepositoryURL)
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

func Initialize() *RepositoryServer {
	return &RepositoryServer{
		Repositories: make(map[string]*gopassRepo),
	}
}

func initializeNewGopassRepository(repositoryUrl string) (*gopassRepo, error) {
	repoDir, err := ioutil.TempDir("", "gopass")
	if err != nil {
		log.Printf("not able to create local repository directory: %v", err)
		return nil, err
	}

	err = cloneGopassRepo(repositoryUrl, repoDir)
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

func cloneGopassRepo(repositoryUrl string, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repositoryUrl,
		Depth:    1,
		Progress: os.Stdout,
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
