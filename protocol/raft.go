package protocol

import "github.com/funkytennisball/hera/common"

// Raft consensus protocol
type Raft struct {
}

// Insert (and replaces) an entry given a key
func (r Raft) Insert(k common.Key, e common.Entity) error {
	return nil
}

// Get an entry given a key
func (r Raft) Get(k common.Key) common.Entity {
	return common.Entity{V: "dummy"}
}

// Remove an entry given a key
func (r Raft) Remove(k common.Key) error {
	return nil
}
