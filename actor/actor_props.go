package actor

import "light-actor-go/pid"

type ActorProps struct {
	parent *pid.PID
}

func NewActorProps(parent *pid.PID) *ActorProps {
	props := new(ActorProps)
	props.parent = parent
	return props
}

func ConfigureActorProps(props ...ActorProps) *ActorProps {
	if len(props) > 0 {
		return &props[0]
	}
	return defaultConfig()
}

func (prop *ActorProps) AddParent(parent *pid.PID) {
	prop.parent = parent
}

func defaultConfig() *ActorProps {
	props := new(ActorProps)
	props.parent = nil
	return props
}
