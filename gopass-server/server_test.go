package gopass_server

import (
	"context"
	"os"
	"path/filepath"
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
