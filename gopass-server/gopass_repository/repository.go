package gopass_repository

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"github.com/mdreem/gopass-operator/gopass-server/gopass_repository/cluster"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	ssh_2 "golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

type config struct {
	Path string `yaml:"path"`
}

func (r *RepositoryServer) initializeRepository(ctx context.Context, repository *gopass_repository.Repository) error {
	log.Printf("InitializeRepository called with: %s", (*repository).RepositoryURL)

	_, ok := (r.Repositories)[(*repository).RepositoryURL]
	if ok {
		log.Printf("repository with URL '%s' already initialized", (*repository).RepositoryURL)
		return nil
	}

	credentials, err := r.Client.GetRepositoryCredentials(ctx, repository.Authentication)
	if err != nil {
		log.Printf("error initializing repository: %v", err)
		return err
	}

	err = r.Client.GetGpgKey(ctx, repository.Authentication)
	if err != nil {
		log.Printf("error fetching gpgKey: %v", err)
		return err
	}

	gopassRepository, err := initializeNewGopassRepository((*repository).RepositoryURL, credentials)
	if err != nil {
		log.Printf("error initializing repository: %v", err)
		return err
	}

	(r.Repositories)[(*repository).RepositoryURL] = gopassRepository

	return nil
}

func (r *RepositoryServer) updateRepository(ctx context.Context, repository *gopass_repository.Repository) error {
	log.Printf("UpdateRepository called with: %s", (*repository).RepositoryURL)

	repo, ok := (r.Repositories)[(*repository).RepositoryURL]
	if !ok {
		log.Printf("unable to find repository with with URL '%s'", (*repository).RepositoryURL)

		return fmt.Errorf("unable to find repository with with URL '%s'", (*repository).RepositoryURL)
	}

	credentials, err := r.Client.GetRepositoryCredentials(ctx, repository.Authentication)
	if err != nil {
		log.Printf("error initializing repository: %v", err)
		return err
	}

	err = updateGopassRepo(repo.repository, credentials.Name, credentials.Password)
	if err != nil {
		return err
	}
	log.Printf("synced repository with URL '%s'", (*repository).RepositoryURL)

	return nil
}

func initializeNewGopassRepository(repositoryUrl string, credentials cluster.Secret) (*gopassRepo, error) {
	repoDir, err := ioutil.TempDir("", "gopass")
	if err != nil {
		log.Printf("not able to create local repository directory: %v", err)
		return nil, err
	}

	repository, err := cloneGopassRepo(repositoryUrl, repoDir, credentials.Name, credentials.Password)
	if err != nil {
		log.Printf("not able clone gopass repository with URL %s: %v", repositoryUrl, err)
		return nil, err
	}

	store, err := createNewGopassClient(context.Background(), repoDir)
	if err != nil {
		log.Printf("not able to create new gopass client: %v", err)
		return nil, err
	}

	gr := &gopassRepo{
		store:      store,
		directory:  repoDir,
		repository: repository,
	}

	return gr, nil
}

func cloneGopassRepo(repositoryUrl string, path string, username string, password string) (*git.Repository, error) {
	sshPassword := ssh.Password{
		User:     username,
		Password: password,
		HostKeyCallbackHelper: ssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh_2.InsecureIgnoreHostKey(),
		},
	}

	log.Printf("cloning repository with URL '%s' to %s\n", repositoryUrl, path)
	repository, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      repositoryUrl,
		Progress: os.Stdout,
		Auth:     &sshPassword,
	})
	return repository, err
}

func createNewGopassClient(ctx context.Context, path string) (*api.Gopass, error) {
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

	return store, nil
}

func removeFile(file *os.File) {
	err := os.Remove(file.Name())
	if err != nil {
		log.Printf("failed to remove file: %v\n", err)
	}
}

func updateGopassRepo(repository *git.Repository, username string, password string) error {
	sshPassword := ssh.Password{
		User:     username,
		Password: password,
		HostKeyCallbackHelper: ssh.HostKeyCallbackHelper{
			HostKeyCallback: ssh_2.InsecureIgnoreHostKey(),
		},
	}

	log.Printf("pulling repository\n")
	worktree, err := repository.Worktree()
	if err != nil {
		log.Printf("unable to fetch worktree of repository: %v\n", err)
		return err
	}
	err = worktree.Pull(&git.PullOptions{
		Auth: &sshPassword,
	})
	if err != nil {
		if err == git.NoErrAlreadyUpToDate {
			log.Printf("repository already up to date")
			return nil
		}
		log.Printf("unable to pull changes: %v", err)
		return err
	}
	return nil
}
