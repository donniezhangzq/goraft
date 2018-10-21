package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"time"
)

type Election struct {
	role constant.ElectionState
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
}

func NewElection(options *Options, logger *log.Logger) *Election {
	return &Election{
		role:               constant.Follower,
		voteFor:            options.Id,
		members:            options.Members,
		ElectionMinTimeout: options.ElectionTimeoutMin,
		ElectionMaxTimeout: options.ElectionTimeoutMax,
		logger:             logger,
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

func (e *Election) heartBeat() {

}

func (e *Election) discoverMembers() {

}

func (e *Election) contest() {
	e.role = constant.Candidate

}
