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
