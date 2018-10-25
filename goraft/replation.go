package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"fmt"
	"net/rpc"
	"sync"
	"time"
)

type Replation struct {
	currentTerm     int
	currentTermMut  *sync.Mutex
	index           int
	commonInfo      *CommonInfo
	entries         *Entries
	entriesToCommit chan *Entry
	lastLogIndex    int
	lastLogTerm     int
	//index of highest log entry to be committed
	commitIndex int
	//index of hithest log entry applied to state machine
	lastApplied int
	//log number to sync at once
	tocommitBufferSize int
	logger             *log.Logger
	rpcClientCache     *RpcClientCache
}

type Entry struct {
	key        interface{}
	value      interface{}
	term       int
	index      int
	isCommited bool
}

type Entries struct {
	entry []*Entry
	mut   *sync.Mutex
}

func NewEntries() *Entries {
	return &Entries{
		entry: make([]*Entry, 0),
		mut:   &sync.Mutex{},
	}
}

func NewReplation(options *Options, logger *log.Logger, rpcClientCache *RpcClientCache, commonInfo *CommonInfo) *Replation {
	return &Replation{
		currentTerm:        0,
		currentTermMut:     &sync.Mutex{},
		index:              0,
		commonInfo:         commonInfo,
		entries:            NewEntries(),
		entriesToCommit:    make(chan *Entry),
		lastLogIndex:       0,
		lastLogTerm:        0,
		commitIndex:        0,
		lastApplied:        0,
		tocommitBufferSize: options.TocommitBufferSize,
		logger:             logger,
		rpcClientCache:     rpcClientCache,
	}
}

func (r *Replation) Start() error {
	r.logger.Debug("replation model start")
	return nil
}

func (r *Replation) Stop() error {
	return nil
}

func (r *Replation) GetCurrentTerm() int {
	r.currentTermMut.Lock()
	defer r.currentTermMut.Unlock()
	return r.currentTerm
}

func (r *Replation) SetCurrentTerm(term int) error {
	r.currentTermMut.Lock()
	defer r.currentTermMut.Unlock()
	if term < r.currentTerm {
		return constant.ErrTermLessThanCurrentTerm
	}
	r.currentTerm = term
	return nil
}

func (r *Replation) GetLastLogIndex() int {
	return r.lastLogIndex
}

func (r *Replation) GetLastLogTerm() int {
	return r.lastLogTerm
}

//send  rpc request for append entries to followers and candidate
func (r *Replation) hearbeat() {
	if r.commonInfo.GetRole() != constant.Leader {
		r.logger.Error(constant.ErrOnlyLeaderCanSendHearbeat.Error())
		return
	}
	//todo async send append entris rpc request
	for id, client := range r.rpcClientCache.GetRpcClients() {
		if id == r.commonInfo.Id {
			continue
		}
		response := new(RpcAppendEntriesResponse)
		r.requestAppendEntries(id, client, response)
		r.logger.Info(fmt.Sprintf("heartbeat response is %v", response))
	}
}

//send rpc request for append entries to followers and candidate with a client
func (r *Replation) requestAppendEntries(id string, client *rpc.Client, response *RpcAppendEntriesResponse) {
	args := &RpcAppendEntriesReqArgs{}
	if err := client.Call(constant.RpcAppendEntries, args, response); err != nil {
		r.logger.Error(fmt.Sprintf("send append entries rpc request failed,target:%s,Error:%s", id, err.Error()))
	}
}

func (r *Replation) AppendEntries(args *RpcAppendEntriesReqArgs, response *RpcAppendEntriesResponse) error {
	//todo complete the appendEntries alghriam
	if args.Term < r.currentTerm {
		return constant.ErrHearbeatTermLess
	}
	r.commonInfo.SetLastHearbeattime(time.Now())
	r.commonInfo.SetLeaderId(args.LeaderId)

	return nil
}
