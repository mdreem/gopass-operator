package gopass_repository

import (
	"archive/zip"
	"context"
	"github.com/go-git/go-git/v5"
	config2 "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/mdreem/gopass-operator/gopass-server/gopass_repository/cluster"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "gopass repository can be initialized",
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			gpgIdFilename := filepath.Join(dir, ".gpg-id")
			gpgIdFile, err := os.Create(gpgIdFilename)
			if err != nil {
				t.Errorf("not able to create temporary file: %s", gpgIdFilename)
				return
			}
			err = gpgIdFile.Close()
			if err != nil {
				t.Errorf("not able close temporary file: %s", gpgIdFilename)
				return
			}

			_, err = createNewGopassClient(tt.args.ctx, dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("createNewGopassClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInitializeRepository(t *testing.T) {
	repoDir := initializeTestRepository(t)

	r := RepositoryServer{
		Repositories: map[string]*gopassRepo{},
		Client:       &cluster.KubernetesTestClient{},
	}

	repositoryInitialization := gopass_repository.RepositoryInitialization{
		Repository: &gopass_repository.Repository{
			RepositoryURL: repoDir,
			Authentication: &gopass_repository.Authentication{
				Namespace: "testNameSpace",
				Username:  "testUsername",
				SecretRef: "testSecretRef",
				SecretKey: "testSecretKey",
			},
			SecretName: nil,
		},
	}

	response, err := r.InitializeRepository(context.Background(), &repositoryInitialization)
	if err != nil {
		t.Errorf("not able to initialize repository: %v\n", err)
		return
	}

	if !response.Successful {
		t.Error("response not successful")
	}

	if response.ErrorMessage != "" {
		t.Errorf("received error message '%s', expected empty error message", response.ErrorMessage)
	}

	deleteDirectory(t, repoDir)
	deleteDirectory(t, r.Repositories[repoDir].directory)
}

func TestCloneRepository(t *testing.T) {
	repoDir := initializeTestRepository(t)
	targetDir := t.TempDir()

	unzip(filepath.Join("resources_test", "password-store.zip"), repoDir, t)

	_, err := cloneGopassRepo(repoDir, targetDir, "", "")
	if err != nil {
		t.Errorf("not able to clone repository: %v", err)
		return
	}
}

func TestInitializeNewGopassRepository(t *testing.T) {
	repoDir := initializeTestRepository(t)

	repository, err := initializeNewGopassRepository(repoDir, cluster.Secret{})
	if err != nil {
		t.Errorf("not able to initialize gopass repository: %v\n", err)
		return
	}

	deleteDirectory(t, repository.directory)
}

func TestCloneAndUpdateRepository(t *testing.T) {
	localRepoDir := initializeTestRepository(t)
	unzip(filepath.Join("resources_test", "password-store.zip"), localRepoDir, t)

	targetDir := t.TempDir()
	repository, err := cloneGopassRepo(localRepoDir, targetDir, "", "")
	if err != nil {
		t.Errorf("unable to clone repository: %v\n", err)
	}

	filename := filepath.Join(localRepoDir, "hello-world-file")
	err = ioutil.WriteFile(filename, []byte("hello world!"), 0644)
	if err != nil {
		t.Errorf("unable to write to file: %v", err)
	}

	localRepo, err := git.PlainOpen(localRepoDir)
	if err != nil {
		t.Errorf("unable to open repository: %v", err)
	}

	worktree, err := localRepo.Worktree()
	if err != nil {
		t.Errorf("unable to get worktree: %v", err)
	}

	_, err = worktree.Add("hello-world-file")
	if err != nil {
		t.Errorf("unable to add file: %v", err)
	}

	_, err = worktree.Commit("some commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	r, repo := createRepositoryServer(repository)
	err = r.updateRepository(context.Background(), &repo)
	if err != nil {
		t.Errorf("unable to update repository: %v\n", err)
	}

}

func TestUpdateGopassRepository(t *testing.T) {
	remoteRepo := initializeTestRepository(t)
	localRepoDir, init, _ := initializeLocalRepository(t)

	_, err := init.CreateRemote(&config2.RemoteConfig{
		Name: "origin",
		URLs: []string{
			remoteRepo,
		},
		Fetch: nil,
	})
	if err != nil {
		t.Errorf("unable to create remote: %v\n", err)
	}

	r, repo := createRepositoryServer(init)
	response, err := r.UpdateRepository(context.Background(), &repo)
	if err != nil {
		t.Errorf("unable to update repository: %v\n", err)
	}

	if !response.Successful {
		t.Error("response not successful")
	}

	if response.ErrorMessage != "" {
		t.Errorf("received error message '%s', expected empty error message", response.ErrorMessage)
	}

	files, err := ioutil.ReadDir(localRepoDir)
	if err != nil {
		t.Errorf("unable to read directory: %v", err)
	}

	foundFiles := make(map[string]bool)
	for _, f := range files {
		foundFiles[f.Name()] = true
		t.Log(f.Name())
	}

	if len(foundFiles) != 3 {
		t.Errorf("expected exactly 3 files, but found %d", len(foundFiles))
	}

	if !foundFiles[".gpg-id"] {
		t.Errorf("expected .git in the pulled repository")
	}

	if !foundFiles[".git"] {
		t.Errorf("expected .git in the pulled repository")
	}

	if !foundFiles["testpwd.gpg"] {
		t.Errorf("expected .git in the pulled repository")
	}

	deleteDirectory(t, remoteRepo)
	deleteDirectory(t, localRepoDir)
}

func createRepositoryServer(init *git.Repository) (RepositoryServer, gopass_repository.Repository) {
	gr := gopassRepo{
		store:      nil,
		directory:  "",
		repository: init,
	}

	r := RepositoryServer{
		Repositories: map[string]*gopassRepo{
			"myRepo": &gr,
		},
		Client: &cluster.KubernetesTestClient{},
	}
	return r, gopass_repository.Repository{
		RepositoryURL:  "myRepo",
		Authentication: nil,
		SecretName:     nil,
	}
}

func initializeLocalRepository(t *testing.T) (string, *git.Repository, error) {
	localRepoDir := t.TempDir()

	init, err := git.PlainInit(localRepoDir, false)
	if err != nil {
		t.Errorf("unable to initialize repository: %v\n", err)
		return "", nil, err
	}
	t.Logf("created repository in %s", localRepoDir)
	return localRepoDir, init, nil
}

func deleteDirectory(t *testing.T, directory string) {
	t.Logf("removing: %s", directory)
	if strings.HasPrefix(directory, os.TempDir()) {
		err := os.RemoveAll(directory)
		if err != nil {
			t.Errorf("not able to remove directory (%s): %v\n", directory, err)
		}
	}
}

func initializeTestRepository(t *testing.T) string {
	repoDir := t.TempDir()

	unzip(filepath.Join("resources_test", "password-store.zip"), repoDir, t)

	return filepath.Join(repoDir, ".password-store")
}

func unzip(src string, dest string, t *testing.T) {
	r, err := zip.OpenReader(src)
	if err != nil {
		t.Errorf("not able to open repository: %v", err)
	}
	defer r.Close()

	for _, f := range r.File {
		t.Logf("extracting: %s\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			t.Errorf("not able to open file (%s): %v", f.Name, err)
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				t.Errorf("not able to create directory (%s): %v", path, err)
			}
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				t.Errorf("not able to open file for writing (%s): %v", path, err)
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				t.Errorf("unable to write to file (%s): %v", path, err)
			}
		}
	}
}
