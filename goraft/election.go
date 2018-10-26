package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"fmt"
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
	electionTimes       int
	contestMut          *sync.Mutex
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
		electionTimes:      0,
		contestMut:         &sync.Mutex{},
	}
}

func (e *Election) start() error {
	e.logger.Debug("enter election start")
	//start a election
	e.contestMut.Lock()
	e.contest()
	e.contestMut.Unlock()

	go e.monitorHearbeatLost()

	e.logger.Debug("election model start")
	return nil
}

func (e *Election) stop() error {
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

func (e *Election) requestForVote(id string, client *rpc.Client, response *RpcElectionResponse) {
	e.logger.Debug(fmt.Sprintf("id:%s, client:%v prepare to send vote request", id, client))
	var args = new(RpcElectionReqArgs)
	args.CandidateId = e.commonInfo.Id
	args.Term = e.replation.getCurrentTerm()
	args.LastLogIndex = e.replation.getLastLogIndex()
	args.LastLogTerm = e.replation.getLastLogTerm()
	e.logger.Debug("args is:", args)
	//todo async call vote
	if err := client.Call(constant.RpcVote, args, response); err != nil {
		e.logger.Warn("request vote from id:%s failed,Error:%s", id, err.Error())
	}
	e.logger.Debug(fmt.Sprintf("complete request vote,response is %v", response))
}

func (e *Election) contest() {
	e.electionTimes += 1
	if e.electionTimes == constant.MaxElectionTimes {
		e.logger.Error("election failed until maxElectionTimes")
		e.electionTimes = 0
		e.commonInfo.transferToFollower()
		return
	}

	e.logger.Debug("election start to contest")

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
		e.electionTimes = 0
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

	//the server hed transfered to leader or follower
	e.electionTimes = 0
}

func (e *Election) monitorHearbeatLost() {
	var ticker = time.NewTicker(constant.MonitorHearbeatlostInterval)
	defer ticker.Stop()

	for {
		r := e.getRandomElecationTimeout()
		timeout := e.commonInfo.lastHeartbeatTime.Add(r)
		if time.Now().Before(timeout) {
			e.contestMut.Lock()
			e.contest()
			e.contestMut.Unlock()
		}

		select {
		case <-ticker.C:
		}
	}
}

//vote for candidate
func (e *Election) Vote(args *RpcElectionReqArgs, response *RpcElectionResponse) error {
	e.logger.Debug("receive vote rpc request.args:%v", args)
	response.Term = args.Term
	response.VoteGranted = false

	e.logger.Debug("getVotefor for is:", e.getVoteFor())
	e.logger.Debug("currentTerm is:", e.replation.currentTerm)

	if args.Term < e.replation.getCurrentTerm() {
		//not grant vote
		e.logger.Debug("not grant vote, term error")
		return nil
	}
	//follower will deny its vote if its own log is more up-do-date
	if args.LastLogTerm < e.replation.getLastLogTerm() || args.LastLogIndex < e.replation.getLastLogIndex() {
		//not grant vote
		e.logger.Debug("not grant vote, log error")
		return nil
	}
	if e.getVoteFor() == "" || e.getVoteFor() == args.CandidateId {
		//grant vote
		e.logger.Debug("grant vote")
		response.VoteGranted = true
		e.setVoteFor(args.CandidateId)
		e.logger.Debug("grant vote response is", response)
	}
	e.logger.Debug("response is", response)
	return nil
}
