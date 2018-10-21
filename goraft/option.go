package goraft

import (
	"donniezhangzq/goraft/constant"
	"donniezhangzq/goraft/utils"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
	"time"
)

type Options struct {
	//master id
	Id string
	//GOMAXPROCS number
	Core int
	//client http request port
	ClintPort int
	//rpc port
	RpcPort int
	//members' address and rpc ports
	Members *Members
	//client http request prefix
	HttpPrefix string
	LogLevel   string
	LogPath    string
	//vote rpc timeout min and max number
	ElectionTimeoutMin time.Duration
	ElectionTimeoutMax time.Duration
	//Candidate not get majority votes timeout to revote
	CandidateTimeoutMin time.Duration
	CandidateTimeoutMax time.Duration
	//to commit buffer size
	TocommitBufferSize int
	//Replication append interval
	ReplationSyncInterval time.Duration
}

func NewOption() *Options {

	return &Options{
		Core:                  2,
		ClintPort:             8299,
		RpcPort:               8199,
		Members:               NewMembers(),
		HttpPrefix:            "goraft",
		LogLevel:              "debug",
		LogPath:               "./goraft.log",
		ElectionTimeoutMin:    time.Duration(150 * time.Millisecond),
		ElectionTimeoutMax:    time.Duration(300 * time.Millisecond),
		CandidateTimeoutMin:   time.Duration(100 * time.Millisecond),
		CandidateTimeoutMax:   time.Duration(200 * time.Millisecond),
		TocommitBufferSize:    1000,
		ReplationSyncInterval: time.Duration(50 * time.Millisecond),
	}
}

func (op *Options) ParseOptions(config *ini.File) error {
	var (
		defaultCfg     = config.Section(constant.DefaultSection)
		logCfg         = config.Section(constant.LogSection)
		electionCfg    = config.Section(constant.ElectionSection)
		replicationCfg = config.Section(constant.ReplicationSection)
	)

	var err error

	if defaultCfg.HasKey("id") {
		op.Id = defaultCfg.Key("id").String()
	} else {
		return constant.ErrIdNotExist
	}
	if defaultCfg.HasKey("core") {
		op.Core, err = defaultCfg.Key("core").Int()
		if err != nil {
			return err
		}
	}
	if defaultCfg.HasKey("client_port") {
		op.ClintPort, err = defaultCfg.Key("client_port").Int()
		if err != nil {
			return err
		}
	}
	if defaultCfg.HasKey("rpc_port") {
		op.RpcPort, err = defaultCfg.Key("rpc_port").Int()
		if err != nil {
			return err
		}
	}
	//to check this if options.Id in members
	var id_include = false
	//parse members to struct member
	if defaultCfg.HasKey("members") {
		members := strings.Split(defaultCfg.Key("members").String(), constant.MembersConfigSep)
		if len(members) < constant.MinMemberNumber {
			return constant.ErrMembersNumberLess
		}
		//member's number must be even number
		if len(members)%2 == 0 {
			return constant.ErrMembersNumberIsNotOdd
		}
		for _, member := range members {
			address := strings.Split(member, constant.AddressPortSep)
			if len(address) != 3 {
				return constant.ErrAddressGetFailed
			}
			id := address[0]
			if id == op.Id {
				id_include = true
			}

			ip := address[1]
			if !utils.CheckIpAddress(ip) {
				return constant.ErrIpformatError
			}
			port, err := strconv.Atoi(address[2])
			if err != nil {
				return err
			}
			m := &Member{
				Id:      id,
				Address: ip,
				RpcPort: port,
				Alive:   true,
				Role:    constant.Follower,
			}
			op.Members.members = append(op.Members.members, m)
		}
	} else {
		return constant.ErrMembersNotExist
	}

	if !id_include {
		return constant.ErrIdNotInMembers
	}

	if defaultCfg.HasKey("http_prefix") {
		op.HttpPrefix = defaultCfg.Key("http_prefix").String()
	}
	if logCfg.HasKey("log_level") {
		op.LogLevel = logCfg.Key("log_level").String()
	}
	if logCfg.HasKey("log_path") {
		op.LogPath = logCfg.Key("log_path").String()
	}
	if electionCfg.HasKey("election_timeout_min") {
		etmin, err := electionCfg.Key("election_timeout_min").Int()
		if err != nil {
			return err
		}
		op.ElectionTimeoutMin = time.Duration(etmin) * time.Millisecond
	}
	if electionCfg.HasKey("election_timeout_max") {
		etmax, err := electionCfg.Key("election_timeout_max").Int()
		if err != nil {
			return err
		}
		op.ElectionTimeoutMax = time.Duration(etmax) * time.Millisecond
	}
	if electionCfg.HasKey("candidate_timeout_min") {
		ctmin, err := electionCfg.Key("candidate_timeout_min").Int()
		if err != nil {
			return err
		}
		op.CandidateTimeoutMin = time.Duration(ctmin) * time.Millisecond
	}
	if electionCfg.HasKey("candidate_timeout_max") {
		ctmax, err := electionCfg.Key("candidate_timeout_max").Int()
		if err != nil {
			return err
		}
		op.CandidateTimeoutMax = time.Duration(ctmax) * time.Millisecond
	}
	if replicationCfg.HasKey("tocommit_buffer_size") {
		op.TocommitBufferSize, err = replicationCfg.Key("tocommit_buffer_size").Int()
		if err != nil {
			return err
		}
	}
	if replicationCfg.HasKey("replation_sync_interval") {
		rsi, err := replicationCfg.Key("replation_sync_interval").Int()
		if err != nil {
			return err
		}
		op.ReplationSyncInterval = time.Duration(rsi) * time.Millisecond
	}
	return nil
}
