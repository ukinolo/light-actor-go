package remote

import (
	context "context"
	"fmt"
	"light-actor-go/actor"
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

type RemoteConfig struct {
	Addr string
}

type RemoteReceiver struct {
	UnimplementedRemoteReceiverServer
	actorSystem *actor.ActorSystem
	server      *grpc.Server
	config      *RemoteConfig
}

func NewRemoteConfig(addr string) *RemoteConfig {
	return &RemoteConfig{Addr: addr}
}

func NewRemoteReceiver(config *RemoteConfig, actorSystem *actor.ActorSystem) *RemoteReceiver {
	receiver := &RemoteReceiver{
		config:      config,
		actorSystem: actorSystem,
	}

	return receiver
}

// startServer starts the gRPC server and listens for incoming messages
func (r *RemoteReceiver) startServer() {
	lis, err := net.Listen("tcp", r.config.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	r.server = grpc.NewServer()
	RegisterRemoteReceiverServer(r.server, r)
	log.Printf("server listening at %v", lis.Addr())
	if err := r.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (r *RemoteReceiver) ReceiveMessage(context.Context, *Envelope) (*Error, error) {
	//TODO find actor
	//TODO send message to actor
	//TODO return
	fmt.Println("Pogodio sam metodu") //TODO remove
	return &Error{Error: ""}, nil
}
