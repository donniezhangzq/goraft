package goraft

import (
	"donniezhangzq/goraft/constant"
	"sync"
)

type Member struct {
	Id      string
	Address string
	RpcPort int
	Alive   bool
	Role    constant.ElectionState
}

func NewMember(id, address string, rpcPort int, role constant.ElectionState) *Member {
	return &Member{
		Id:      id,
		Address: address,
		RpcPort: rpcPort,
		Alive:   true,
		Role:    role,
	}
}

type Members struct {
	members []*Member
	mut     *sync.Mutex
}

func NewMembers() *Members {
	return &Members{
		members: make([]*Member, 0),
		mut:     &sync.Mutex{},
	}
}

func (m *Members) checkIdUniq() bool {
	for i := 0; i < len(m.members)-1; i++ {
		for j := i + 1; j < len(m.members); j++ {
			if m.members[i] == m.members[j] {
				return false
			}
		}
	}
	return true
}

func (m *Members) GetMembers() []*Member {
	return m.members
}
