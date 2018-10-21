package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"donniezhangzq/goraft/storage"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Server interface {
	Start() error
	Stop() error
}

type Goraft struct {
	Id         string
	ClientPort int
	RpcPort int
	HttpPrefix string
	election   *Election
	replation  *Replation
	storage    *storage.Storage
	isStart    bool
	logger     *log.Logger
}

func NewGoraft(options *Options, logger *log.Logger) (*Goraft, error) {
	election, err := NewElection(options, logger)
	if err != nil {
		return nil, err
	}
	storage := storage.NewStorage(election)
	replation := NewReplation(options, logger)

	return &Goraft{
		Id:         options.Id,
		RpcPort:options.RpcPort,
		ClientPort: options.ClintPort,
		HttpPrefix: options.HttpPrefix,
		election:   election,
		replation:  replation,
		storage:    storage,
		isStart:    false,
		logger:     logger,
	}, nil
}

func (g *Goraft) Start() error {
	if err := g.replation.Start(); err != nil {
		return err
	}

	if err := g.election.Start(); err != nil {
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
