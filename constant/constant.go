package constant

import "time"

//leader election state
type ElectionState int

const (
	Leader ElectionState = iota
	Candidate
	Follower
)

const (
	DefaultConfigPath = "./goraft.conf"
)

const (
	MinMemberNumber  = 3
	MembersConfigSep = ","
	AddressPortSep   = ":"
)

//configuration section
const (
	DefaultSection     = "default"
	LogSection         = "log"
	ElectionSection    = "election"
	ReplicationSection = "replication"
)

//hashmap count
const (
	BucketCount = 16
)

//utils.checkip.regexp
const (
	MatchIp = `^([0-9]{1,3}\.){3}[0-9]{1,3}$`
)

//rpc method
const (
	RpcVote          = "Election.Vote"
	RpcAppendEntries = "Replation.AppendEntries"
)

//election
const (
	MaxElectionTimes            = 5
	MonitorHearbeatlostInterval = 50 * time.Microsecond
)
