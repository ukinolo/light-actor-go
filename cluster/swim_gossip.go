package cluster

import "fmt"

// type SwimState int32

// const (
// 	Alive SwimState = iota
// 	Suspicious
// 	Dead
// )

// type SwimPing struct {
// 	sender string
// 	extras SwimExtraInfo
// }

// type SwimAck struct {
// 	sender  string
// 	state   SwimState
// 	version int
// 	extras  SwimExtraInfo
// }

// type SwimIndirectPing struct {
// 	receiver     string
// 	intermediate string
// 	sender       string
// }

// type SwimIndirectAck struct {
// 	receiver string
// 	sender   string
// 	state    SwimState
// 	version  int
// }

type SendPing struct{} //TODO make private
type JoinCluster struct {
	address string
}

// type SwimExtraInfo struct {
// 	member_address string
// 	member_state   SwimState
// 	version        int
// 	ttl            int //Time to live
// }

// type SwimNewGossiper struct {
// 	member_address string
// 	member_state   SwimState
// 	version        int
// }

type RespondWaitInfo struct {
	waitTime              int
	indirectPingRequested bool //I sent indirect ping request to someone about this node
}

// SWIM: Scalable Weakly-consistent Infection-style Process Group Membership Protocol
type SwimGossiper struct {
	memberList          *MemberList
	healthyAddresses    []string
	noResponseAddresses map[string]RespondWaitInfo //List of addresses that we pinged, and did not received ack
	deadInterval        int                        //How many ping intervals I will wait for ack before I announce someone dead
	nextToPing          int
	extraInfo           *SwimExtraInfo
	selfAddress         string
	selfVersion         int
}

func NewSwimGossiper(memberList *MemberList, selfAddress string) *SwimGossiper {
	addresses := make([]string, 0)
	for k, v := range memberList.mapping {
		if v.state == SwimState_Alive {
			addresses = append(addresses, k)
		}
	}
	return &SwimGossiper{
		memberList:          memberList,
		healthyAddresses:    addresses,
		nextToPing:          0,
		noResponseAddresses: make(map[string]RespondWaitInfo),
		selfAddress:         selfAddress,
		selfVersion:         1,
		deadInterval:        5, //TODO make parameter
	}
}

func (SwimGossiper *SwimGossiper) HandlePing(ping *SwimPing) *SwimAck {
	ack := SwimAck{
		Sender:  SwimGossiper.selfAddress,
		State:   SwimState_Alive,
		Version: int32(SwimGossiper.selfVersion),
		Extras:  SwimGossiper.extraInfo,
	}

	SwimGossiper.handleExtraInfo(ping.Extras)
	return &ack
}

func (SwimGossiper *SwimGossiper) HandleAck(ack *SwimAck) {
	SwimGossiper.memberList.Update(ack.Sender, ClusterInfo{state: ack.State, version: int(ack.Version)})
	//TODO change this
	//This is necesary because I can ping someone who is not in my member list if he pinged me
	//So I need to update healthyAddresses list if he is not in there
	if ack.GetState() == SwimState_Alive {
		found := false
		for _, v := range SwimGossiper.healthyAddresses {
			if v == ack.Sender {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Pomoglo je 0000000000000000000000000000000000000000000000000000000000000000000")
			SwimGossiper.healthyAddresses = append(SwimGossiper.healthyAddresses, ack.Sender)
		}
	}

	delete(SwimGossiper.noResponseAddresses, ack.Sender)

	SwimGossiper.handleExtraInfo(ack.Extras)

	if SwimGossiper.extraInfo != nil && SwimGossiper.extraInfo.Ttl <= 1 {
		SwimGossiper.extraInfo = &SwimExtraInfo{
			MemberAddress: ack.Sender,
			MemberState:   ack.State,
			Version:       ack.Version,
			Ttl:           3, //TODO make parameter
		}
	}
}

func (SwimGossiper *SwimGossiper) HandleIndirectPing(ping *SwimIndirectPing) *SwimIndirectAck {
	ack := SwimIndirectAck{
		Receiver: ping.Sender,
		Sender:   SwimGossiper.selfAddress,
		State:    SwimState_Alive,
		Version:  int32(SwimGossiper.selfVersion),
	}

	return &ack
}

func (SwimGossiper *SwimGossiper) HandleIndirectAck(ack *SwimIndirectAck) {
	ci := SwimGossiper.memberList.Find(ack.Sender)
	if ci.state == SwimState_Dead && ack.State == SwimState_Alive {
		SwimGossiper.healthyAddresses = append(SwimGossiper.healthyAddresses, ack.Sender)
	}
	SwimGossiper.memberList.Update(ack.Sender, ClusterInfo{
		state:   ack.State,
		version: int(ack.Version),
	})
	delete(SwimGossiper.noResponseAddresses, ack.Sender)
}

func (SwimGossiper *SwimGossiper) CreateNewPing() (*SwimPing, *SwimIndirectPing, string) {
	SwimGossiper.selfVersion += 1
	address := SwimGossiper.getNewPingAddress()
	indirectPingAddress := SwimGossiper.updateNoResponseTime(address != "")
	if address == "" {
		return &SwimPing{}, &SwimIndirectPing{}, ""
	}
	ping := SwimPing{
		Sender: SwimGossiper.selfAddress,
		Extras: SwimGossiper.extraInfo,
	}
	indirectPing := SwimIndirectPing{}
	if indirectPingAddress != "" {
		indirectPing.Receiver = indirectPingAddress
		indirectPing.Sender = SwimGossiper.selfAddress
		indirectPing.Intermediate = address
	} else {
		SwimGossiper.noResponseAddresses[address] = RespondWaitInfo{
			waitTime:              0,
			indirectPingRequested: false}
	}

	if SwimGossiper.extraInfo != nil {
		SwimGossiper.extraInfo.Ttl -= 1
	}

	return &ping, &indirectPing, address
}

func (SwimGossiper *SwimGossiper) HandleNewGossiper(newGossiper *SwimNewGossiper) {
	SwimGossiper.memberList.Add(newGossiper.MemberAddress, ClusterInfo{
		state:   newGossiper.MemberState,
		version: int(newGossiper.Version),
	})
	SwimGossiper.healthyAddresses = append(SwimGossiper.healthyAddresses, newGossiper.MemberAddress)
	SwimGossiper.extraInfo = &SwimExtraInfo{
		MemberAddress: newGossiper.MemberAddress,
		MemberState:   newGossiper.MemberState,
		Version:       newGossiper.Version,
		Ttl:           6, //TODO change to parameter
	}
}

// Returns next address to ping
func (SwimGossiper *SwimGossiper) getNewPingAddress() string {
	l := len(SwimGossiper.healthyAddresses)
	if l == 0 {
		return ""
	}
	for range SwimGossiper.healthyAddresses {
		if SwimGossiper.nextToPing >= l-1 {
			SwimGossiper.nextToPing = 0
		} else {
			SwimGossiper.nextToPing += 1
		}
		_, ok := SwimGossiper.noResponseAddresses[SwimGossiper.healthyAddresses[SwimGossiper.nextToPing]]
		if !ok {
			return SwimGossiper.healthyAddresses[SwimGossiper.nextToPing]
		}
	}
	return ""
}

func (SwimGossiper *SwimGossiper) updateNoResponseTime(isAddressForIndirectAvailable bool) string {
	indirectPingCandadate := ""
	var toBeRemoved []string
	for k, v := range SwimGossiper.noResponseAddresses {
		if v.waitTime >= SwimGossiper.deadInterval {
			if v.indirectPingRequested {
				ci := SwimGossiper.memberList.Find(k)

				SwimGossiper.memberList.Update(k, ClusterInfo{state: SwimState_Dead, version: ci.version + 1})
				toBeRemoved = append(toBeRemoved, k)
				continue
			}
			if isAddressForIndirectAvailable {
				if indirectPingCandadate == "" {
					indirectPingCandadate = k
					SwimGossiper.noResponseAddresses[k] = RespondWaitInfo{waitTime: 0,
						indirectPingRequested: true}
				}
			} else {
				ci := SwimGossiper.memberList.Find(k)

				SwimGossiper.memberList.Update(k, ClusterInfo{state: SwimState_Dead, version: ci.version + 1})
				toBeRemoved = append(toBeRemoved, k)
				continue
			}
		}
	}
	for _, k := range toBeRemoved {
		for i, v := range SwimGossiper.healthyAddresses {
			if v == k {
				SwimGossiper.healthyAddresses = append(SwimGossiper.healthyAddresses[:i], SwimGossiper.healthyAddresses[i+1:]...)
				break
			}
		}
		delete(SwimGossiper.noResponseAddresses, k)
	}
	for k := range SwimGossiper.noResponseAddresses { //Map is not the best data structure to do this because of this for loop
		SwimGossiper.noResponseAddresses[k] = RespondWaitInfo{
			waitTime:              SwimGossiper.noResponseAddresses[k].waitTime + 1,
			indirectPingRequested: SwimGossiper.noResponseAddresses[k].indirectPingRequested,
		}
	}
	return indirectPingCandadate
}

func (SwimGossiper *SwimGossiper) handleExtraInfo(extra *SwimExtraInfo) {
	if extra == nil {
		return
	}
	if extra.MemberAddress == "" || extra.Ttl <= 0 || extra.MemberAddress == SwimGossiper.selfAddress {
		return
	}

	ci := SwimGossiper.memberList.Find(extra.MemberAddress)
	if ci.version > int(extra.Version) {
		return
	}
	SwimGossiper.memberList.Update(extra.MemberAddress, ClusterInfo{
		state:   extra.MemberState,
		version: int(extra.Version),
	})
	if ci.version == 0 || extra.GetMemberState() != ci.state {
		if extra.GetMemberState() == SwimState_Alive {
			delete(SwimGossiper.noResponseAddresses, extra.MemberAddress)
			SwimGossiper.healthyAddresses = append(SwimGossiper.healthyAddresses, extra.MemberAddress)
		} else if extra.GetMemberState() == SwimState_Dead {
			for i, v := range SwimGossiper.healthyAddresses {
				if v == extra.MemberAddress {
					SwimGossiper.healthyAddresses = append(SwimGossiper.healthyAddresses[:i], SwimGossiper.healthyAddresses[i+1:]...)
				}
			}
			delete(SwimGossiper.noResponseAddresses, extra.MemberAddress)
		}
	}
}
