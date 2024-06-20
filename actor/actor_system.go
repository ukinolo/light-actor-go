package actor

import (
	"context"
	"light-actor-go/remote"
)

type ActorSystem struct {
	remoteHandler *remote.RemoteHandler
	registry      *Registry
}

// Creates new actor system that can only be used localy
func NewActorSystem() *ActorSystem {
	return &ActorSystem{remoteHandler: nil, registry: NewRegistry()}
}

// Creates new actor system that can be used remotely
func NewRemoteActorSystem(remoteConfig remote.RemoteConfig) *ActorSystem {
	return &ActorSystem{remoteHandler: remote.NewRemoteHandler(&remoteConfig), registry: NewRegistry()}
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
