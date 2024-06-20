package remote

import (
	"context"
	"fmt"
	"light-actor-go/envelope"
	"light-actor-go/pid"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type RemoteConfig struct {
	Addr string
}

type RemoteHandler struct {
	UnimplementedRemoteHandlerServer
	MessageChan  chan *MessageWrapper
	EnvelopeChan *chan envelope.Envelope // Updated to pointer to channel
	server       *grpc.Server
	config       *RemoteConfig
}

func Configure(addr string) *RemoteConfig {
	return &RemoteConfig{Addr: addr}
}

func NewRemoteHandler(config *RemoteConfig, envelopeChan *chan envelope.Envelope) *RemoteHandler { // Updated to pointer to channel
	handler := &RemoteHandler{
		MessageChan:  make(chan *MessageWrapper, 100),
		EnvelopeChan: envelopeChan,
		config:       config,
	}

	go handler.startServer()
	return handler
}

// startServer starts the gRPC server and listens for incoming messages
func (r *RemoteHandler) startServer() {
	lis, err := net.Listen("tcp", r.config.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	r.server = grpc.NewServer()
	RegisterRemoteHandlerServer(r.server, r)
	log.Printf("server listening at %v", lis.Addr())
	if err := r.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// SendMessage handles incoming messages
func (r *RemoteHandler) SendMessage(ctx context.Context, wrapper *MessageWrapper) (*Empty, error) {
	r.MessageChan <- wrapper
	fmt.Printf("Queued message for receiver: %s\n", wrapper.Receiver)
	return &Empty{}, nil
}

// ReceiveMessage streams messages to a receiver
func (r *RemoteHandler) ReceiveMessage(req *MessageRequest, stream RemoteHandler_ReceiveMessageServer) error {
	receiverUUID, err := uuid.Parse(req.Receiver)
	if err != nil {
		return err
	}

	// Receives messages from queued messages to be processed
	for msg := range r.MessageChan {
		msgUUID, err := uuid.Parse(msg.Receiver)
		if err != nil {
			return err
		}

		//Ensures that only messages intended for a specific client (based on the receiver UUID) are sent to that client through the gRPC stream
		if msgUUID == receiverUUID {
			fmt.Printf("Sending message to receiver: %s\n", msgUUID)
			envelope := envelope.NewEnvelope(msg.Message, pid.PID{ID: receiverUUID})
			*r.EnvelopeChan <- *envelope // Updated to send the value through the dereferenced pointer to the channel
			if err := stream.Send(msg); err != nil {
				return err
			}
		}
	}
	return nil
}

// Sends a message to the specified address
func (r *RemoteHandler) SendMessageToAddress(addr string, msg proto.Message, receiverUUID uuid.UUID) error {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials())) // Corrected method name to grpc.Dial
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewRemoteHandlerClient(conn)

	anyMsg, err := anypb.New(msg)
	if err != nil {
		return err
	}

	wrapper := &MessageWrapper{
		Message:  anyMsg,
		Receiver: receiverUUID.String(),
	}

	_, err = client.SendMessage(context.Background(), wrapper)
	return err
}
