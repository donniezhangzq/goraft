package goraft

import (
	"donniezhangzq/goraft/log"
)

type Replation struct {
	term         int
	index        int
	members      *Members
	Log          []*Log
	LogToCommit  chan *Log
	lastLogIndex int
	lastLogTerm  int
	//index of highest log entry to be committed
	commitIndex int
	//index of hithest log entry applied to state machine
	lastApplied int
	//log number to sync at once
	tocommitBufferSize int
	logger             *log.Logger
	rpcClientCache     *RpcClientCache
}

type Log struct {
	key        interface{}
	value      interface{}
	term       int
	index      int
	isCommited bool
}

func NewReplation(options *Options, logger *log.Logger, rpcClientCache *RpcClientCache) *Replation {
	return &Replation{
		term:               1,
		index:              0,
		members:            options.Members,
		Log:                make([]*Log, 0),
		LogToCommit:        make(chan *Log),
		lastLogIndex:       0,
		lastLogTerm:        1,
		commitIndex:        0,
		lastApplied:        0,
		tocommitBufferSize: options.TocommitBufferSize,
		logger:             logger,
		rpcClientCache:     rpcClientCache,
	}
}

func (r *Replation) Start() error {

}

func (r *Replation) Stop() error {

}
