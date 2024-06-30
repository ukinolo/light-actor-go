package actor

type ActorProps struct {
	Parent *PID
}

func NewActorProps(Parent *PID) *ActorProps {
	props := new(ActorProps)
	props.Parent = Parent
	return props
}

func ConfigureActorProps(props ...ActorProps) *ActorProps {
	if len(props) > 0 {
		return &props[0]
	}
	return defaultConfig()
}

func (prop *ActorProps) AddParent(Parent *PID) {
	prop.Parent = Parent
}

func defaultConfig() *ActorProps {
	props := new(ActorProps)
	props.Parent = nil
	return props
}
