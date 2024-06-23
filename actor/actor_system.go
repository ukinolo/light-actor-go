package actor

import (
	"context"
)

type ActorSystem struct {
	registry *Registry
}

// Creates new actor system that can only be used localy
func NewActorSystem() *ActorSystem {
	return &ActorSystem{registry: NewRegistry()}
}

func (system *ActorSystem) SpawnActor(a Actor, props ...ActorProps) (PID, error) {
	prop := ConfigureActorProps(props...)

	actorChan := make(chan Envelope)
	mailbox := NewMailbox(actorChan)

	mailboxChan := mailbox.GetChan()
	mailboxPID, err := NewPID()
	if err != nil {
		return mailboxPID, err
	}

	//Start mailbox in separate gorutine
	StartWorker(mailbox.Start, nil)

	//Start actor in separate gorutine
	StartWorker(func() {
		//Setup basic actor context
		actorContext := NewActorContext(context.Background(), system, prop, mailboxPID)
		for {
			envelope := <-actorChan
			//Set only message and send
			actorContext.AddEnvelope(envelope)
			a.Receive(*actorContext)
		}
	}, nil)

	//Put mailbox chanel in registry
	err = system.registry.Add(mailboxPID, mailboxChan)
	if err != nil {
		return mailboxPID, err
	}

	return mailboxPID, nil
}

func (system *ActorSystem) Send(envelope Envelope) {
	ch := system.registry.Find(*envelope.Receiver())
	if ch == nil {
		return
	}
	ch <- envelope
}

func (system *ActorSystem) AddRemoteActor(remoteActorPID PID, senderChan chan Envelope) {
	system.registry.Add(remoteActorPID, senderChan)
}
