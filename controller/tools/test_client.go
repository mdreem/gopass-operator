package main

import (
	"context"
	"fmt"
	"github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"google.golang.org/grpc"
	"log"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("not able to connect: %v", err)
	}
	defer conn.Close()

	c := gopass_repository.NewRepositoryServiceClient(conn)

	repository, err := c.InitializeRepository(
		context.Background(),
		&gopass_repository.Repository{
			RepositoryURL: "ssh://localhost/home/git/password-store",
		},
	)
	if repository == nil {
		fmt.Printf("got an empty response")
		return
	}
	fmt.Printf("repository: %t", (*repository).Successful)
}
