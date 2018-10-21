package goraft

import (
	"fmt"
	"net/rpc"
	"sync"
)

type RpcClientCache struct {
	rpcClientCache map[string]*rpc.Client
	mut            *sync.Mutex
}

func NewRpcClientCache() *RpcClientCache {
	return &RpcClientCache{
		rpcClientCache: make(map[string]*rpc.Client),
		mut:            &sync.Mutex{},
	}
}

func (r *RpcClientCache) AddRpcClient(member *Member) error {
	r.mut.Lock()
	defer r.mut.Unlock()
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%d", member.Address, member.RpcPort))
	if err != nil {
		return err
	}
	r.rpcClientCache[member.Id] = client
	return nil
}

func (r *RpcClientCache) GetRpcClients() map[string]*rpc.Client {
	return r.rpcClientCache
}

type RpcElectionReqArgs struct {
}

type RpcElectionResponse struct {
}
