package protocol

type Key struct {
	K string
}

type Entity struct {
	V string
}

type Protocol interface {
	Get(k Key) Entity
	Insert(k Key, e Entity)
	Remove(k Key)
}
