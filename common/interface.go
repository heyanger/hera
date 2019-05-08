package common

type Key struct {
	K string
}

type Entity struct {
	V string
}

type Node struct {
	Source    string
	Heartbeat uint64
	Location  string
	Leader    bool
}

type RangeKey struct {
	Start, End   string
	RingLocation uint64
}

type Range struct {
	Nodes []string
	Size  uint64
}

type RangeMap map[RangeKey]Range
type NodeMap map[string]Node
type Store map[string]string

type ServersConfig struct {
	Servers map[string][]string
	// Clientport   string
	// Protocolport string
}

type Protocol interface {
	Get(k Key) Entity
	Insert(k Key, e Entity) error
	Remove(k Key) error
	State() string
	Location() string

	Init(id string, bind string, genesis bool) error
	Join(nodeid string, addr string) error
}
