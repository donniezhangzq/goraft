package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/log"
	"fmt"
	"sync"
	"time"
)

type CommonInfo struct {
	Id                string
	role              constant.ElectionState
	roleMut           *sync.Mutex
	Members           *Members
	logger            *log.Logger
	leaderId          string
	lastHeartbeatTime time.Time
}

func NewCommonInfo(options *Options, logger *log.Logger) *CommonInfo {
	return &CommonInfo{
		Id:                options.Id,
		role:              constant.Follower,
		roleMut:           &sync.Mutex{},
		Members:           NewMembers(),
		logger:            logger,
		lastHeartbeatTime: time.Now(),
	}
}

func (c *CommonInfo) GetRole() constant.ElectionState {
	return c.role
}

func (c *CommonInfo) setRole(role constant.ElectionState) {
	c.role = role
}

func (c *CommonInfo) transferToLeader() {
	c.roleMut.Lock()
	defer c.roleMut.Unlock()
	beforeRole := c.GetRoleName()

	c.setRole(constant.Leader)
	c.SetLeaderId(c.Id)
	c.logger.SetDefaultField(c.GetRoleName(), c.Id, c.GetLeaderId())

	c.logger.Info(fmt.Sprintf("%s transfer to leader", beforeRole))
}

func (c *CommonInfo) transferToCandidate() {
	c.roleMut.Lock()
	defer c.roleMut.Unlock()
	beforeRole := c.GetRoleName()

	c.setRole(constant.Candidate)
	c.logger.SetDefaultField(c.GetRoleName(), c.Id, c.GetLeaderId())

	c.logger.Info(fmt.Sprintf("%s transfer to candidate", beforeRole))
}

func (c *CommonInfo) transferToFollower() {
	c.roleMut.Lock()
	defer c.roleMut.Unlock()
	beforeRole := c.GetRoleName()

	c.setRole(constant.Follower)
	c.logger.SetDefaultField(c.GetRoleName(), c.Id, c.GetLeaderId())

	c.logger.Info(fmt.Sprintf("%s transfer to follower", beforeRole))
}

func (c *CommonInfo) GetLeaderId() string {
	return c.leaderId
}

func (c *CommonInfo) SetLeaderId(id string) {
	c.leaderId = id
}

func (c *CommonInfo) GetLastHeartbeatTime() time.Time {
	return c.lastHeartbeatTime
}

func (c *CommonInfo) SetLastHearbeattime(t time.Time) {
	c.lastHeartbeatTime = t
}

func (c *CommonInfo) GetRoleName() string {
	role := c.GetRole()
	switch role {
	case constant.Leader:
		return "leader"
	case constant.Candidate:
		return "candidate"
	case constant.Follower:
		return "follower"
	default:
		return ""
	}
}
