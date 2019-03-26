package protocol

type Raft struct {
}

func (r Raft) Insert(k Key, e Entity) bool {
	return true
}

func (r Raft) Get(k Key) bool {
	return false
}

func (r Raft) Remove(k Key) bool {
	return false
}
