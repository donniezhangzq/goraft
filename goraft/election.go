package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"math/rand"
	"net/rpc"
	"sync"
	"time"
)

type Election struct {
	commonInfo *CommonInfo
	//id who the server vote for
	voteFor       string
	voteForMut    *sync.Mutex
	leaderAddress string
	//follower to candidate timeout(random range from min to max)
	ElectionMinTimeout time.Duration
	ElectionMaxTimeout time.Duration
	//candidate revote timeout(random range from min to max)
	CandidateMinTimeout time.Duration
	CandidateMaxTimeout time.Duration
	logger              *log.Logger
	rpcClientCache      *RpcClientCache
	replation           *Replation
}

func NewElection(options *Options, logger *log.Logger, rpcclientCache *RpcClientCache, replation *Replation,
	commonInfo *CommonInfo) *Election {

	return &Election{
		commonInfo:         commonInfo,
		voteFor:            options.Id,
		voteForMut:         &sync.Mutex{},
		ElectionMinTimeout: options.ElectionTimeoutMin,
		ElectionMaxTimeout: options.ElectionTimeoutMax,
		logger:             logger,
		rpcClientCache:     rpcclientCache,
		replation:          replation,
	}
}

func (e *Election) Start() error {
	//start a election
	e.contest()
	return nil
}

func (e *Election) Stop() error {
	return nil
}

//get id the server vote for
func (e *Election) getVoteFor() string {
	e.voteForMut.Lock()
	defer e.voteForMut.Unlock()
	return e.voteFor
}

func (e *Election) setVoteFor(id string) {
	e.voteForMut.Lock()
	defer e.voteForMut.Unlock()
	e.voteFor = id
}

func (e *Election) getRandomElecationTimeout() time.Duration {
	return time.Duration(e.getRandMinToMax(int64(e.ElectionMinTimeout), int64(e.ElectionMaxTimeout)))
}

func (e *Election) getRandomCandidateTimeout() time.Duration {
	return time.Duration(e.getRandMinToMax(int64(e.CandidateMinTimeout), int64(e.CandidateMaxTimeout)))
}

func (e *Election) getRandMinToMax(min int64, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	randf := rand.Float64()
	return min + int64(float64(max-min)*randf)
}

//vote for candidate
func (e *Election) Vote(args *RpcElectionReqArgs, response *RpcElectionResponse) error {
	response.Term = args.Term
	response.VoteGranted = false
	if args.Term < e.replation.GetCurrentTerm() {
		//not grant vote
		return nil
	}
	//follower will deny its vote if its own log is more up-do-date
	if args.LastLogTerm < e.replation.GetLastLogTerm() || args.LastLogIndex < e.replation.GetLastLogIndex() {
		//not grant vote
		return nil
	}
	if e.getVoteFor() == "" || e.getVoteFor() == args.CandidateId {
		//grant vote
		response.VoteGranted = true
		e.setVoteFor(args.CandidateId)
	} else {
		//not grant vote
		return nil
	}
	return nil
}

func (e *Election) requestForVote(id string, client *rpc.Client, response *RpcElectionResponse) {
	var args = new(RpcElectionReqArgs)
	args.CandidateId = e.commonInfo.Id
	args.Term = e.replation.GetCurrentTerm()
	args.LastLogIndex = e.replation.GetLastLogIndex()
	args.LastLogTerm = e.replation.GetLastLogTerm()
	//todo async call vote
	if err := client.Call(constant.RpcVote, args, response); err != nil {
		e.logger.Warn("request vote from id:%s failed,Error:%s", id, err.Error())
	}
}

func (e *Election) contest() {
	//first step:transfer follower to candidate
	e.commonInfo.transferToCandidate()
	//second step:request for vote
	//todo async call vote
	var voteCount = 0
	var winVote = false
	for id, client := range e.rpcClientCache.GetRpcClients() {
		if id != e.commonInfo.Id {
			response := new(RpcElectionResponse)
			e.requestForVote(id, client, response)
			if response.VoteGranted {
				//grant vote
				voteCount += 1
			}
			if voteCount > int(e.commonInfo.Members.GetCount()/2) {
				winVote = true
				break
			}
		}
	}
	//third step:transfer to leader(win) or remain candidate(fail)
	if winVote {
		//win the vote and send heartbeat to others
		e.commonInfo.transferToLeader()
		e.replation.hearbeat()
	} else {
		//remain candidate and wait a time to receive new leader's heartbeat
		time.Sleep(e.getRandomCandidateTimeout())
		if e.commonInfo.GetRole() == constant.Candidate {
			//try to request vote again
			e.contest()
		}
	}
}
