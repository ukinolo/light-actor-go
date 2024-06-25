package actor

type systemMessage interface {
	getAction() string
}

type ActorStarted struct{}

type ActorStoped struct{}

//System messages

type forcefulShutdown struct{}

func (forcefulShutdown) getAction() string {
	return "Forceful shutdown"
}

type gracefulShutdown struct{}

func (gracefulShutdown) getAction() string {
	return "Greceful shutdown"
}

type startingActor struct{}

func (startingActor) getAction() string {
	return "Starting actor"
}

type stopingActor struct{}

func (stopingActor) getAction() string {
	return "Stopping actor"
}
