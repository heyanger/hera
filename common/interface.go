package common

type NodeKey uint32

type Key struct {
	K string
}

type Entity struct {
	V string
}

type Node struct {
	Source    string
	Heartbeat uint64
}

type RangeKey struct {
	Start, End string
}

type Range struct {
	Nodes      [3]NodeKey
	RaftLeader NodeKey
	Size       uint64
}

type ServersConfig struct {
	Servers map[string][]string
	// Clientport   string
	// Protocolport string
}

type Protocol interface {
	Get(k Key) Entity
	Insert(k Key, e Entity) error
	Remove(k Key) error

	Init(nodes []string, localID string) error
}
