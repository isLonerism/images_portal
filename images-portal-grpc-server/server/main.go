package main

import (
	"fmt"
	"log"
	"net"

	"github.com/bsuro10/images_portal/images-portal-grpc-server/api/docker"
	"github.com/docker/docker/client"
	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	dockerService := docker.RegisterDockerService(dockerClient)

	grpcServer := grpc.NewServer()

	docker.RegisterDockerServer(grpcServer, dockerService)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
