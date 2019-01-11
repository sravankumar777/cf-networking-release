package fakes

import (
	"bosh-dns-adapter/copilot/api"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type CopilotVIPResolverServer struct {
	listener net.Listener
	server   *grpc.Server
}

func (c *CopilotVIPResolverServer) Start(port int) {
	var err error
	c.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	c.server = grpc.NewServer()
	api.RegisterVIPResolverCopilotServer(c.server, c)

	go func() {
		c.server.Serve(c.listener)
	}()
}

func (c *CopilotVIPResolverServer) Close() error {
	return c.listener.Close()
}

func (c *CopilotVIPResolverServer) Health(context.Context, *api.HealthRequest) (*api.HealthResponse, error) {
	return &api.HealthResponse{Healthy: true}, nil
}

func (c *CopilotVIPResolverServer) GetVIPByName(ctx context.Context, request *api.GetVIPByNameRequest) (*api.GetVIPByNameResponse, error) {
	return &api.GetVIPByNameResponse{}, nil
}
