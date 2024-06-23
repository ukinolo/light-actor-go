package remote

import (
	"context"

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

func (rs *RemoteSender) SendMessage(message interface{}, receiverName string) error {
	conn, err := grpc.NewClient(rs.remoteAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	client := NewRemoteReceiverClient(conn)

	anyMsg, err := anypb.New(message.(proto.Message))
	if err != nil {
		return err
	}

	protoEnvelope := &Envelope{
		Message:  anyMsg,
		Receiver: receiverName,
	}

	_, err = client.ReceiveMessage(context.Background(), protoEnvelope)
	if err != nil {
		return err
	}
	return nil
}
