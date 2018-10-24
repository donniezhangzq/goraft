package constant

import (
	"errors"
	"fmt"
)

//storage error
var (
	ErrNodeIsNil       = errors.New("node is nil")
	ErrNodeConflictKey = errors.New("conflict key")
	ErrRoleIsNotLeader = errors.New("key can be changed only by leader")
)

//utils
var (
	ErrGetLocalIpFailed = errors.New("get local ip failed")
)

//optons
var (
	ErrIdNotExist            = errors.New("id not exist in config file")
	ErrMinValueMoreThanMax   = errors.New("min value is not less than max value")
	ErrMembersNotExist       = errors.New("members not exist in config file")
	ErrMembersNumberLess     = errors.New(fmt.Sprintf("members' number less thean expect:%d", MinMemberNumber))
	ErrMembersNumberIsNotOdd = errors.New("members's number is not odd")
	ErrAddressGetFailed      = errors.New("get address from member failed")
	ErrIpformatError         = errors.New("ip address format error")
	ErrIdNotInMembers        = errors.New("id config not in members config")
	ErrIdInMembersDuplicate  = errors.New("id in members is not uniq")
)

//replation
var (
	ErrTermLessThanCurrentTerm = errors.New("term in less than currentTerm")
)

//election
var (
	ErrOnlyLeaderCanSendHearbeat = errors.New("only leader can send heartbeat")
	ErrHearbeatTermLess          = errors.New("hearbeat args's term is less then the server's term")
)

//log
var (
	ErrLogPathIsNotFile = errors.New("log path is not a file")
)
