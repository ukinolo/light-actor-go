package cluster

import (
	"sync"
)

type ClusterInfo struct {
	state   SwimState
	version int
}

type MemberList struct {
	mapping map[string]ClusterInfo
	sm      sync.RWMutex
}

func NewMemberList() *MemberList {
	return &MemberList{mapping: make(map[string]ClusterInfo)}
}

func (memberList *MemberList) Find(address string) ClusterInfo {
	memberList.sm.RLock()
	defer memberList.sm.RUnlock()
	return memberList.mapping[address]
}

// Add does not look at version
func (memberList *MemberList) Add(address string, ci ClusterInfo) {
	memberList.sm.Lock()
	defer memberList.sm.Unlock()
	memberList.mapping[address] = ci
}

func (memberList *MemberList) Delete(address string, ci ClusterInfo) {
	memberList.sm.Lock()
	defer memberList.sm.Unlock()
	delete(memberList.mapping, address)
}

// Update checks for version
func (memberList *MemberList) Update(address string, ci ClusterInfo) {
	memberList.sm.Lock()
	defer memberList.sm.Unlock()
	current, ok := memberList.mapping[address]
	if !ok || memberList.mapping[address].version < current.version {
		memberList.mapping[address] = ci
	}
}

func (memberList *MemberList) GetAllAlive() []string {
	aliveAddresses := make([]string, 0)
	memberList.sm.RLock()
	defer memberList.sm.RUnlock()
	for k, v := range memberList.mapping {
		if v.state == SwimState_Alive {
			aliveAddresses = append(aliveAddresses, k)
		}
	}
	return aliveAddresses
}
