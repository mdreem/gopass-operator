package gopass_server

import (
	"archive/zip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestCloneRepository(t *testing.T) {
	repoDir := t.TempDir()
	targetDir := t.TempDir()

	unzip(filepath.Join("resources_test", "password-store.zip"), repoDir, t)

	err := cloneGopassRepo(filepath.Join(repoDir, ".password-store"), targetDir)
	if err != nil {
		t.Errorf("not able to clone repository: %v", err)
		return
	}
}

func TestInitializeNewGopassRepository(t *testing.T) {
	repoDir := t.TempDir()

	unzip(filepath.Join("resources_test", "password-store.zip"), repoDir, t)

	repository, err := InitializeNewGopassRepository(filepath.Join(repoDir, ".password-store"))
	if err != nil {
		t.Errorf("not able to initialize gopass repository: %v\n", err)
		return
	}

	t.Logf("removing: %s", repository.directory)
	if strings.HasPrefix(repository.directory, os.TempDir()) {
		err = os.RemoveAll(repository.directory)
		if err != nil {
			t.Errorf("not able to remove directory (%s): %v\n", repository.directory, err)
			return
		}
	}
}

func unzip(src string, dest string, t *testing.T) {
	r, err := zip.OpenReader(src)
	if err != nil {
		t.Errorf("not able to open repository: %v", err)
		return
	}
	defer r.Close()

	for _, f := range r.File {
		t.Logf("extracting: %s\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			t.Errorf("not able to open file (%s): %v", f.Name, err)
			return
		}
		defer rc.Close()

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, f.Mode())
			if err != nil {
				t.Errorf("not able to create directory (%s): %v", path, err)
				return
			}
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				t.Errorf("not able to open file for writing (%s): %v", path, err)
				return
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				t.Errorf("unable to write to file (%s): %v", path, err)
				return
			}
		}
	}
}
