package goraft

import (
	"donniezhangzq/goraft/log"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"net/rpc"
)

type Server interface {
	Start() error
	Stop() error
}

type Goraft struct {
	ClientPort     int
	RpcPort        int
	HttpPrefix     string
	election       *Election
	replation      *Replation
	storage        *Storage
	isStart        bool
	logger         *log.Logger
	members        *Members
	rpcClientCache *RpcClientCache
}

func NewGoraft(options *Options, logger *log.Logger, rpcClientCache *RpcClientCache) *Goraft {
	commonInfo := NewCommonInfo(options, logger)
	replation := NewReplation(options, logger, rpcClientCache, commonInfo)
	election := NewElection(options, logger, rpcClientCache, replation, commonInfo)
	storage := NewStorage(commonInfo)

	return &Goraft{
		RpcPort:        options.RpcPort,
		ClientPort:     options.ClintPort,
		HttpPrefix:     options.HttpPrefix,
		election:       election,
		replation:      replation,
		storage:        storage,
		isStart:        false,
		logger:         logger,
		members:        options.Members,
		rpcClientCache: rpcClientCache,
	}
}

func (g *Goraft) Start() error {
	if err := g.startRpc(); err != nil {
		return err
	}

	for _, member := range g.members.GetMembers() {
		if err := g.rpcClientCache.AddRpcClient(member); err != nil {
			g.logger.Fatal(fmt.Sprintf("addRpcClientCache failed,member:%v,Error:%s", member, err.Error()))
		}
	}

	if err := g.election.Start(); err != nil {
		return err
	}

	if err := g.replation.Start(); err != nil {
		return err
	}

	g.isStart = true

	g.logger.Debug("goraft start")

	if err := g.startHttp(); err != nil {
		return err
	}
	return nil
}

func (g *Goraft) Stop() error {
	if !g.isStart {
		return nil
	}

	if err := g.replation.Stop(); err != nil {
		return err
	}

	if err := g.election.Stop(); err != nil {
		return err
	}
	return nil
}

func (g *Goraft) GetKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (g *Goraft) PostKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (g *Goraft) DelKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (g *Goraft) List(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (g *Goraft) startHttp() error {
	router := httprouter.New()
	storagePath := fmt.Sprintf("/%s/storage/:name", g.HttpPrefix)
	router.GET(storagePath, g.GetKey)
	router.POST(storagePath, g.PostKey)
	router.DELETE(storagePath, g.DelKey)
	router.GET(fmt.Sprintf("/%s/storage", g.HttpPrefix), g.List)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", g.ClientPort), router); err != nil {
		g.logger.Fatal("start clint http server failed,Error", err.Error())
		return err
	}
	return nil
}

func (g *Goraft) startRpc() error {
	rpc.Register(g.election)
	rpc.Register(g.replation)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", g.RpcPort))
	if err != nil {
		g.logger.Fatal("listen error:%s", err.Error())
	}
	go http.Serve(l, nil)
	return nil
}
