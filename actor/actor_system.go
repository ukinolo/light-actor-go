package actor

import (
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

func (system *ActorSystem) SpawnActor(a Actor) (PID, error) {
	actorChan := make(chan Envelope)
	mailbox := NewMailbox(actorChan)

	//Start mailbox in separate gorutine
	StartWorker(mailbox.Start, nil)

	//Start actor in separate gorutine
	StartWorker(func() {
		//Setup basic actor context
		for {
			<-actorChan
			//Set only message and send
			a.Recieve(ActorContext{})
		}
	}, nil)

	mailboxChan := mailbox.GetChan()
	mailboxPID, err := NewPID()
	if err != nil {
		return mailboxPID, err
	}

	//Put mailbox chanel in registry
	err = system.registry.Add(mailboxPID, mailboxChan)
	if err != nil {
		return mailboxPID, err
	}

	return mailboxPID, err
}

func (system *ActorSystem) Send(envelope Envelope, reciever PID) {
	ch := system.registry.Find(reciever)
	if ch == nil {
		return
	}
	ch <- envelope
}
