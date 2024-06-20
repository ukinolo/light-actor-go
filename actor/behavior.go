package actor

type ReceiveFunc func(context ActorContext)

// Behavior represents a stack of ReceiveFunc functions.
type Behavior struct {
	stack []ReceiveFunc
}

func NewBehavior() *Behavior {
	return &Behavior{
		stack: make([]ReceiveFunc, 0),
	}
}

func (b *Behavior) Become(receive ReceiveFunc) {
	b.clear()
	b.push(receive)
}

func (b *Behavior) BecomeStacked(receive ReceiveFunc) {
	b.push(receive)
}

func (b *Behavior) UnbecomeStacked() {
	b.pop()
}

func (b *Behavior) Receive(context ActorContext) {
	if behavior, ok := b.peek(); ok {
		behavior(context)
	} else {

	}
}

func (b *Behavior) clear() {
	b.stack = b.stack[:0]
}

func (b *Behavior) push(receive ReceiveFunc) {
	b.stack = append(b.stack, receive)
}

func (b *Behavior) pop() (ReceiveFunc, bool) {
	if len(b.stack) == 0 {
		return nil, false
	}

	index := len(b.stack) - 1
	receive := b.stack[index]
	b.stack = b.stack[:index]
	return receive, true
}

func (b *Behavior) peek() (ReceiveFunc, bool) {
	if len(b.stack) == 0 {
		return nil, false
	}

	return b.stack[len(b.stack)-1], true
}

func (b *Behavior) len() int {
	return len(b.stack)
}
