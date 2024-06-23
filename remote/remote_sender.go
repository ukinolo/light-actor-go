package remote

import (
	"context"
	"errors"
	"light-actor-go/actor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

type RemoteSender struct {
	remoteAddress string
}

func NewRemoteSender(address string) *RemoteSender {
	return &RemoteSender{
		remoteAddress: address,
	}
}

func (rs *RemoteSender) SendMessage(envelope actor.Envelope) error {
	conn, err := grpc.NewClient(rs.remoteAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := NewRemoteReceiverClient(conn)

	anyMsg, err := anypb.New(envelope.Message.(proto.Message))
	if err != nil {
		return err
	}

	protoEnvelope := &Envelope{
		Message:  anyMsg,
		Receiver: "Neko",
	}

	ret, err := client.ReceiveMessage(context.Background(), protoEnvelope)
	if ret.Error != "" {
		return errors.New(ret.Error)
	}
	if err != nil {
		return err
	}
	return nil
}
