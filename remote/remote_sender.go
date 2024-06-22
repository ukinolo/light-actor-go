package remote

import (
	"light-actor-go/envelope"
	"light-actor-go/pid"
	"log"

	"google.golang.org/protobuf/proto"
)

type RemoteSender struct {
	remoteHandler *RemoteHandler
	remoteAddress string
	pid           pid.PID
	chanSender    chan envelope.Envelope
}

func NewRemoteSender(handler *RemoteHandler, address string, pid pid.PID) *RemoteSender {
	return &RemoteSender{
		remoteHandler: handler,
		remoteAddress: address,
		pid:           pid,
		chanSender:    make(chan envelope.Envelope),
	}
}

func (rs *RemoteSender) GetChan() chan envelope.Envelope {
	return rs.chanSender
}

func (rs *RemoteSender) Start() {
	for {
		env := <-rs.chanSender
		msg, ok := env.Message.(proto.Message)
		if !ok {
			log.Printf("Failed to send message: not a valid protobuf message")
			continue
		}
		receiverUUID := env.Receiver().ID
		err := rs.remoteHandler.SendMessageToAddress(rs.remoteAddress, msg, receiverUUID)
		if err != nil {
			log.Printf("Failed to send message to %s: %v", rs.remoteAddress, err)
		}
	}
}
