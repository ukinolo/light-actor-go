package remote

import (
	"light-actor-go/actor"
	"log"

	"google.golang.org/protobuf/proto"
)

type RemoteSender struct {
	remoteHandler *RemoteHandler
	remoteAddress string
	pid           actor.PID
	chanSender    chan actor.Envelope
}

func NewRemoteSender(handler *RemoteHandler, address string, pid actor.PID) *RemoteSender {
	return &RemoteSender{
		remoteHandler: handler,
		remoteAddress: address,
		pid:           pid,
		chanSender:    make(chan actor.Envelope),
	}
}

func (rs *RemoteSender) GetChan() chan actor.Envelope {
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
