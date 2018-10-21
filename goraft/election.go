package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"net/rpc"
	"sync"
	"time"
)

type Election struct {
	id      string
	role    constant.ElectionState
	roleMut *sync.Mutex
	//id who the server vote for
	voteFor       string
	members       *Members
	leaderAddress string
	//follower to candidate timeout(random range from min to max)
	ElectionMinTimeout time.Duration
	ElectionMaxTimeout time.Duration
	//candidate revote timeout(random range from min to max)
	CandidateMinTimeout time.Duration
	CandidateMaxTimeout time.Duration
	logger              *log.Logger
	rpcClientCache      *RpcClientCache
}

func NewElection(options *Options, logger *log.Logger, rpcclientCache *RpcClientCache) *Election {
	return &Election{
		id:                 options.Id,
		role:               constant.Follower,
		roleMut:            &sync.Mutex{},
		voteFor:            options.Id,
		members:            options.Members,
		ElectionMinTimeout: options.ElectionTimeoutMin,
		ElectionMaxTimeout: options.ElectionTimeoutMax,
		logger:             logger,
		rpcClientCache:     rpcclientCache,
	}
}

func (e *Election) Start() error {
	//start a election
	e.contest()
}

func (e *Election) Stop() error {

}

func (e *Election) GetRole() constant.ElectionState {
	return e.role
}

func (e *Election) HeartBeat() {

}

//vote for candidate
func (e *Election) Vote() {

}

func (e *Election) requestForVote(id string, client *rpc.Client, response *RpcElectionResponse) {
	var args = new(RpcElectionReqArgs)
	if err := client.Call(constant.RpcVote, args, response); err != nil {
		e.logger.Warn("request vote from id:%s failed,Error:%s", id, err.Error())
	}
}

func (e *Election) contest() {
	//transfer follower to candidate
	e.transterToCandidate()
	//request for vote
	for id, client := range e.rpcClientCache.GetRpcClients() {
		if id != e.id {
			response := new(RpcElectionResponse)
			e.requestForVote(id, client, response)
		}
	}
}

func (e *Election) transterToCandidate() {
	e.role = constant.Candidate
}
