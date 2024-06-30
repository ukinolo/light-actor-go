package cluster

import (
	"light-actor-go/actor"
	"sync"
)

type actorInfo struct {
	address string
	name    string
	created bool
}

type ConnectionActorsRegistry struct {
	mapping map[actor.PID]actorInfo
	sm      sync.RWMutex
}

func NewConnectionActorsRegistry() *ConnectionActorsRegistry {
	return &ConnectionActorsRegistry{mapping: make(map[actor.PID]actorInfo)}
}

func (ConnectionActorsRegistry *ConnectionActorsRegistry) Find(pid actor.PID) actorInfo {
	ConnectionActorsRegistry.sm.RLock()
	defer ConnectionActorsRegistry.sm.RUnlock()
	return ConnectionActorsRegistry.mapping[pid]
}

func (ConnectionActorsRegistry *ConnectionActorsRegistry) Add(pid actor.PID, info actorInfo) {
	ConnectionActorsRegistry.sm.Lock()
	defer ConnectionActorsRegistry.sm.Unlock()
	ConnectionActorsRegistry.mapping[pid] = info
}

func (ConnectionActorsRegistry *ConnectionActorsRegistry) Delete(pid actor.PID) {
	ConnectionActorsRegistry.sm.Lock()
	defer ConnectionActorsRegistry.sm.Unlock()
	delete(ConnectionActorsRegistry.mapping, pid)
}

func (ConnectionActorsRegistry *ConnectionActorsRegistry) GetPid(name string) (actor.PID, bool) {
	ConnectionActorsRegistry.sm.RLock()
	defer ConnectionActorsRegistry.sm.RUnlock()
	for i, v := range ConnectionActorsRegistry.mapping {
		if v.name == name && !v.created {
			return i, true
		}
	}
	return actor.PID{}, false
}
