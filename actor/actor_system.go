package actor

import (
	"context"
	"light-actor-go/envelope"
	"light-actor-go/pid"
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
	remoteServerChan := make(chan envelope.Envelope)
	system := &ActorSystem{remoteHandler: remote.NewRemoteHandler(&remoteConfig, &remoteServerChan), registry: NewRegistry()}
	system.listenRemoteServer(remoteServerChan)
	return system
}

func (system *ActorSystem) listenRemoteServer(remoteServerChan chan envelope.Envelope) {
	for {
		env := <-remoteServerChan
		system.Send(env)
	}
}
func (system *ActorSystem) SpawnActor(a Actor, props ...ActorProps) (pid.PID, error) {
	prop := ConfigureActorProps(props...)

	actorChan := make(chan envelope.Envelope)
	mailbox := NewMailbox(actorChan)

	mailboxChan := mailbox.GetChan()
	mailboxPID, err := pid.NewPID()
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
			a.Recieve(*actorContext)
		}
	}, nil)

	//Put mailbox chanel in registry
	err = system.registry.Add(mailboxPID, mailboxChan)
	if err != nil {
		return mailboxPID, err
	}

	return mailboxPID, nil
}

func (system *ActorSystem) Send(envelope envelope.Envelope) {
	ch := system.registry.Find(*envelope.Receiver())
	if ch == nil {
		return
	}
	ch <- envelope
}

func (system *ActorSystem) SpawnRemoteActor(a Actor, remoteAddress string, props ...ActorProps) (pid.PID, error) {
	prop := ConfigureActorProps(props...)

	remoteSenderPID, err := pid.NewPID()
	if err != nil {
		return remoteSenderPID, err
	}

	// Create RemoteSender
	remoteSender := remote.NewRemoteSender(system.remoteHandler, remoteAddress, remoteSenderPID)

	// Start actor in separate goroutine
	StartWorker(func() {
		// Setup basic actor context
		actorContext := NewActorContext(context.Background(), system, prop, remoteSenderPID)
		for {
			envelope := <-remoteSender.GetChan()
			// Set only message and send
			actorContext.AddEnvelope(envelope)
			a.Recieve(*actorContext)
		}
	}, nil)

	// Put RemoteSender in registry
	err = system.registry.Add(remoteSenderPID, remoteSender.GetChan())
	if err != nil {
		return remoteSenderPID, err
	}

	return remoteSenderPID, nil
}
