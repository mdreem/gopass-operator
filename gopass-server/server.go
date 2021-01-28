package gopass_server

import (
	"context"
	"fmt"
	"github.com/gopasspw/gopass/pkg/gopass/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type gopassRepo struct {
	store *api.Gopass
}

type config struct {
	Path string `yaml:"path"`
}

func createNewGopassClient(ctx context.Context, path string) (*gopassRepo, error) {
	file, err := ioutil.TempFile("", "*config.yml")
	if err != nil {
		return nil, err
	}
	defer funcName(file)
	fmt.Printf("created temporary configuration file: %s\n", file.Name())

	c := config{
		Path: path,
	}

	marshalledConfig, err := yaml.Marshal(&c)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(marshalledConfig)
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	err = os.Setenv("GOPASS_CONFIG", file.Name())
	if err != nil {
		return nil, err
	}

	store, err := api.New(ctx)
	if err != nil {
		return nil, err
	}

	gr := &gopassRepo{
		store: store,
	}

	return gr, nil
}

func funcName(file *os.File) {
	err := os.Remove(file.Name())
	fmt.Printf("failed to remove file: %v", err)
}
