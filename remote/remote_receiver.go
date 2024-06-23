package remote

import (
	context "context"
	"errors"
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
	actorSystem        *actor.ActorSystem
	server             *grpc.Server
	config             *RemoteConfig
	localActorRegistry Registry //Registy of local actors that are discoverable remotely
}

func NewRemoteConfig(addr string) *RemoteConfig {
	return &RemoteConfig{Addr: addr}
}

func NewRemoteReceiver(config *RemoteConfig, actorSystem *actor.ActorSystem) *RemoteReceiver {
	receiver := &RemoteReceiver{
		config:             config,
		actorSystem:        actorSystem,
		localActorRegistry: *NewRegistry(),
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

func (r *RemoteReceiver) AddRemoteActor(name string, actorPID actor.PID) error {
	return r.localActorRegistry.Add(name, actorPID)
}

func (r *RemoteReceiver) ReceiveMessage(context context.Context, envelope *Envelope) (*Empty, error) {
	actorPID := r.localActorRegistry.Find(envelope.Receiver)
	if (actorPID == actor.PID{}) {
		return &Empty{}, errors.New("no actor with name " + envelope.Receiver + " exists")
	}
	actorEnvelope := actor.NewEnvelope(envelope.Message, actorPID)
	r.actorSystem.Send(actorEnvelope)
	return &Empty{}, nil
}
