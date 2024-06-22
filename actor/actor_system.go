package actor

type ActorSystem interface {
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
