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
