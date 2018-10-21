package goraft

import (
	"donniezhangzq/goraft/log"
	"donniezhangzq/goraft/storage"
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
	Id         string
	ClientPort int
	RpcPort    int
	HttpPrefix string
	election   *Election
	replation  *Replation
	storage    *storage.Storage
	isStart    bool
	logger     *log.Logger
}

func NewGoraft(options *Options, logger *log.Logger, rpcClientCache *RpcClientCache) *Goraft {
	election := NewElection(options, logger, rpcClientCache)
	s := storage.NewStorage(election)
	replation := NewReplation(options, logger, rpcClientCache)

	return &Goraft{
		Id:         options.Id,
		RpcPort:    options.RpcPort,
		ClientPort: options.ClintPort,
		HttpPrefix: options.HttpPrefix,
		election:   election,
		replation:  replation,
		storage:    s,
		isStart:    false,
		logger:     logger,
	}
}

func (g *Goraft) Start() error {
	if err := g.election.Start(); err != nil {
		return err
	}
	if err := g.replation.Start(); err != nil {
		return err
	}

	if err := g.startHttp(); err != nil {
		return err
	}
	return nil
}

func (g *Goraft) Stop() error {
	if err := g.replation.Stop(); err != nil {
		return err
	}

	if err := g.election.Stop(); err != nil {
		return err
	}
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
