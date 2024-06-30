package cluster

import (
	"context"
	"fmt"
	"light-actor-go/remote"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type ClusterSender struct{}

func NewClusterSender() *ClusterSender {
	return &ClusterSender{}
}

func (clusterSender *ClusterSender) SendMessage(message interface{}, address string, actorName string) error {
	fmt.Printf("Saljem %v sa %T na %v aktoru %v\n", message, message, address, actorName)
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := remote.NewRemoteReceiverClient(conn)

	anyMsg, err := anypb.New(message.(proto.Message))
	if err != nil {
		return err
	}

	protoEnvelope := &remote.Envelope{
		Message:  anyMsg,
		Receiver: actorName,
	}

	_, err = client.ReceiveMessage(context.Background(), protoEnvelope)
	if err != nil {
		return err
	}
	return nil
}
