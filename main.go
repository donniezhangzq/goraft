package main

import (
	"flag"
	"gopkg.in/ini.v1"
	"errors"
	"fmt"
	"github.com/judwhite/go-svc/svc"
)

type program struct{
	logger *Logger
	goraft *Goraft
}

func initConfig() (*ini.File, error) {
	var configPath string
	flag.StringVar(&configPath, "config", defaultConfigPath, "config file")
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

	options := NewOption()
	if err := options.ParseOptions(config); err != nil {
		panic(fmt.Sprintf("ParseOptions failed,Error:%s", err.Error()))
	}

	logger := NewLogger()
	if err := logger.InitLogger(options); err != nil {
		panic(fmt.Sprintf("init logger failed,Error:%s", err.Error()))
	}

	goraft, err := NewGoraft(options, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("create goraft failed,Error:%s", err.Error()))
		panic(err.Error())
	}

	prg := &program{logger: logger, goraft:goraft}
	if err := svc.Run(prg); err != nil {
		panic(fmt.Sprintf("goraft start failed,Error:%s", err.Error()))
	}
}


func (p *program) Init(env svc.Environment) error {
	p.logger.Debug(fmt.Sprintf("is win service? %v\n", env.IsWindowsService()))
	return nil
}

func (p *program) Start() error {
	p.logger.Debug("goraft startting")
	defer p.logger.Debug("goraft startted")
	return p.goraft.Start()
}

func (p *program) Stop() error {
	p.logger.Debug("goraft stopping")
	defer p.logger.Debug("goraft stopped")
	return p.goraft.Stop()
}