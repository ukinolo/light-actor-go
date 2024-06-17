package actor

type ActorProps struct {
	parent *PID
}

func NewActorProps(parent *PID) *ActorProps {
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

func defaultConfig() *ActorProps {
	props := new(ActorProps)
	props.parent = nil
	return props
}
