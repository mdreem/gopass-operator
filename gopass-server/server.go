package gopass_server

import (
	"fmt"
	"github.com/mdreem/gopass-operator/gopass-server/gopass_repository"
	gopass_repository_grpc "github.com/mdreem/gopass-operator/pkg/apiclient/gopass_repository"
	"google.golang.org/grpc"
	"log"
	"net"
)

func Run() {
	grpcServer := grpc.NewServer()

	gopassRepoServer, err := gopass_repository.Initialize()
	if err != nil {
		log.Fatalf("failed to initialize: %v", err)
	}

	gopass_repository_grpc.RegisterRepositoryServiceServer(grpcServer, gopassRepoServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 9000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to start grpc server: %v", err)
	}
}
