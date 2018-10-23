package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"fmt"
	"sync"
)

type CommonInfo struct {
	Id       string
	role     constant.ElectionState
	roleMut  *sync.Mutex
	Members  *Members
	logger   *log.Logger
	leaderId string
}

func NewCommonInfo(options *Options, logger *log.Logger) *CommonInfo {
	return &CommonInfo{
		Id:      options.Id,
		role:    constant.Follower,
		roleMut: &sync.Mutex{},
		Members: NewMembers(),
		logger:  logger,
	}
}

func (c *CommonInfo) GetRole() constant.ElectionState {
	c.roleMut.Lock()
	defer c.roleMut.Unlock()
	return c.role
}

func (c *CommonInfo) transferToLeader() {
	c.roleMut.Lock()
	defer c.roleMut.Unlock()
	c.role = constant.Leader
	c.SetLeaderId(c.Id)
	c.logger.SetDefaultField(constant.Leader, c.Id, c.GetLeaderId())
	c.logger.Info(fmt.Sprintf("%s transfer to leader", c.GetRole()))
}

func (c *CommonInfo) transferToCandidate() {
	c.roleMut.Lock()
	defer c.roleMut.Unlock()
	c.role = constant.Candidate
	c.logger.SetDefaultField(constant.Candidate, c.Id, c.GetLeaderId())
	c.logger.Info(fmt.Sprintf("%s transfer to candidate", c.GetRole()))
}

func (c *CommonInfo) transferToFollower() {
	c.roleMut.Lock()
	defer c.roleMut.Unlock()
	c.role = constant.Follower
	c.logger.SetDefaultField(constant.Follower, c.Id, c.GetLeaderId())
	c.logger.Info(fmt.Sprintf("%s transfer to follower", c.GetRole()))
}

func (c *CommonInfo) GetLeaderId() string {
	return c.leaderId
}

func (c *CommonInfo) SetLeaderId(id string) {
	c.leaderId = id
}
