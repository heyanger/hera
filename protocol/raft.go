package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	"github.com/funkytennisball/hera/common"

	"github.com/hashicorp/raft"
)

const (
	retainSnapshotCount = 2
	raftTimeout         = 10 * time.Second
)

type command struct {
	Op    string        `json:"op,omitempty"`
	Key   common.Key    `json:"key,omitempty"`
	Value common.Entity `json:"value,omitempty"`
}

// Raft consensus protocol
type Raft struct {
	RaftDir  string
	RaftBind string

	mutex sync.Mutex
	store map[common.Key]common.Entity

	raft   *raft.Raft
	logger *log.Logger
}

// Init inits the store. If enableSingle is set, and there are no existing peers,
// then this node becomes the first node, and therefore leader, of the cluster.
// localID should be the server identifier for this node.
func (r *Raft) Init(nodes []string, localID string) error {
	r.store = make(map[common.Key]common.Entity)
	// Setup Raft configuration.
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)

	fmt.Println(r.RaftBind)

	// Setup Raft communication.
	addr, err := net.ResolveTCPAddr("tcp", r.RaftBind)
	if err != nil {
		return err
	}
	transport, err := raft.NewTCPTransport(r.RaftBind, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return err
	}

	// Create the snapshot store. This allows the Raft to truncate the log.
	snapshots, err := raft.NewFileSnapshotStore(r.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	// Create the log store and stable store.
	var logStore raft.LogStore
	var stableStore raft.StableStore

	logStore = raft.NewInmemStore()
	stableStore = raft.NewInmemStore()

	// Instantiate the Raft systems.
	ra, err := raft.NewRaft(config, r, logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	r.raft = ra

	// if enableSingle {
	configuration := raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}
	ra.BootstrapCluster(configuration)
	// }

	// c := make([]raft.Server, 0, len(nodes))
	// for _, n := range nodes {
	// 	c = append(c, raft.Server{Suffrage: raft.Voter, ID: raft.ServerID(n), Address: raft.ServerAddress(n)})
	// }

	// configuration := raft.Configuration{c}
	// ra.BootstrapCluster(configuration)

	return nil
}

// Get returns the value for the given key.
func (r Raft) Get(key common.Key) common.Entity {
	// s.mu.Lock()
	// defer s.mu.Unlock()
	return r.store[key]
}

// Insert sets the value for the given key.
func (r Raft) Insert(key common.Key, entity common.Entity) error {
	if r.raft.State() != raft.Leader {
		return new(NotALeaderError)
	}

	c := &command{
		Op:    "set",
		Key:   key,
		Value: entity,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := r.raft.Apply(b, raftTimeout)
	return f.Error()
}

// Remove deletes the given key.
func (r Raft) Remove(key common.Key) error {
	if r.raft.State() != raft.Leader {
		return new(NotALeaderError)
	}

	c := &command{
		Op:  "delete",
		Key: key,
	}
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}

	f := r.raft.Apply(b, raftTimeout)
	return f.Error()
}

// Apply applies a Raft log entry to the key-value store.
func (r *Raft) Apply(l *raft.Log) interface{} {
	var c command
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	switch c.Op {
	case "set":
		return r.applySet(c.Key, c.Value)
	case "delete":
		return r.applyDelete(c.Key)
	default:
		panic(fmt.Sprintf("unrecognized command op: %s", c.Op))
	}
}

// Snapshot returns a snapshot of the key-value store.
func (r *Raft) Snapshot() (raft.FSMSnapshot, error) {
	r.mutex.Lock()

	// Clone the map.
	o := make(map[common.Key]common.Entity)
	for k, v := range r.store {
		o[k] = v
	}

	d := &fsmSnapshot{store: o}
	r.mutex.Unlock()

	return d, nil
}

// Restore stores the key-value store to a previous state.
func (r *Raft) Restore(rc io.ReadCloser) error {
	o := make(map[common.Key]common.Entity)
	if err := json.NewDecoder(rc).Decode(&o); err != nil {
		return err
	}

	// Set the state from the snapshot, no lock required according to
	// Hashicorp docs.
	r.store = o
	return nil
}

func (r *Raft) applySet(key common.Key, value common.Entity) interface{} {
	r.mutex.Lock()
	r.store[key] = value
	r.mutex.Unlock()
	return nil
}

func (r *Raft) applyDelete(key common.Key) interface{} {
	r.mutex.Lock()
	delete(r.store, key)
	r.mutex.Unlock()
	return nil
}

type fsmSnapshot struct {
	store map[common.Key]common.Entity
}

func (f *fsmSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		// Encode data.
		b, err := json.Marshal(f.store)
		if err != nil {
			return err
		}

		// Write data to sink.
		if _, err := sink.Write(b); err != nil {
			return err
		}

		// Close the sink.
		return sink.Close()
	}()

	if err != nil {
		sink.Cancel()
	}

	return err
}

func (f *fsmSnapshot) Release() {}
