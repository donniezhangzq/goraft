package main

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/goraft"
	"donniezhangzq/goraft/log"
	"errors"
	"flag"
	"fmt"
	"github.com/judwhite/go-svc/svc"
	logr "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type program struct {
	logger *log.Logger
	goraft *goraft.Goraft
}

func initConfig() (*ini.File, error) {
	var configPath string
	flag.StringVar(&configPath, "config", constant.DefaultConfigPath, "config file")
	flag.Parse()
	if configPath == "" {
		return nil, errors.New("configPath is empty")
	}

	config, err := ini.Load(configPath)
	return config, err
}

func main() {
	config, err := initConfig()
	if err != nil {
		panic(fmt.Sprintf("init config failed,Error:%s", err.Error()))
	}

	options := goraft.NewOption()
	if err := options.ParseOptions(config); err != nil {
		panic(fmt.Sprintf("ParseOptions failed,Error:%s", err.Error()))
	}

	logger := log.NewLogger()
	if err := logger.InitLogger(options.LogPath, options.LogLevel); err != nil {
		panic(fmt.Sprintf("init logger failed,Error:%s", err.Error()))
	}
	//init rpc client cache
	clientCache := goraft.NewRpcClientCache()

	g := goraft.NewGoraft(options, logger, clientCache)

	prg := &program{logger: logger, goraft: g}
	if err := svc.Run(prg); err != nil {
		panic(fmt.Sprintf("goraft start failed,Error:%s", err.Error()))
	}
}

func (p *program) Init(env svc.Environment) error {
	return nil
}

func (p *program) Start() error {
	f := log.NewFatalHook(p.FatalHook, p.logger)
	f.AddHook(f)
	return p.goraft.Start()
}

func (p *program) Stop() error {
	p.logger.Debug("goraft stopping")
	defer p.logger.Debug("goraft stopped")
	return p.goraft.Stop()
}

func (p *program) FatalHook(entry *logr.Entry) error {
	p.logger.Info("enter fatal log hook")
	return p.Stop()
}
