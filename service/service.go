package service

import (
	"encoding/binary"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/funkytennisball/hera/common"
)

// Service provides an HTTP service
type Service struct {
	Protocol common.Protocol

	store  map[common.Key]common.Entity
	ranges common.RangeMap
	nodes  common.NodeMap
}

func (s *Service) postHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: error handling
	w.Write([]byte(s.Protocol.Get(common.Key{K: m["key"]}).V))
}

func (s *Service) putHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: error handling
	if err := s.Protocol.Insert(common.Key{K: m["key"]}, common.Entity{V: m["value"]}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: error handling
	if err := s.Protocol.Remove(common.Key{K: m["key"]}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Upon a new data node, this is
func (s *Service) heartbeat(w http.ResponseWriter, r *http.Request) {
	// TODO: Malicious heartbeat handling
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	keystr := m["key"]
	i, err := strconv.Atoi(keystr)

	if err != nil {
		// Node does not exist, add new entry
		key := common.NodeKey(rand.Uint64())
		v := common.Node{
			Source:    m["source"],
			Heartbeat: uint64(time.Now().UnixNano() / 1000),
		}
		s.nodes[key] = v

		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(key))
		w.Write(b)
		return
	}

	key := common.NodeKey(i)

	if _, ok := s.nodes[key]; ok {
		// Node exists, update heartbeat
		v := s.nodes[key]
		v.Heartbeat = uint64(time.Now().UnixNano() / 1000)
		s.nodes[key] = v
	} else {
		// Node does not exist, add new entry
		key := common.NodeKey(rand.Uint64())
		v := common.Node{
			Source:    m["source"],
			Heartbeat: uint64(time.Now().UnixNano() / 1000),
		}
		s.nodes[key] = v
	}

	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(key))
	w.Write(b)
	return
}

func (s *Service) defaultHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.postHandler(w, r)
	case http.MethodPut:
		s.putHandler(w, r)
	case http.MethodDelete:
		s.deleteHandler(w, r)
	}
}

func (s *Service) heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.heartbeat(w, r)
	}
}

// Start the instance and the HTTP web server
func (s *Service) Start(port string) {
	s.store = make(map[common.Key]common.Entity)
	s.nodes = make(common.NodeMap)
	s.ranges = make(common.RangeMap)

	http.HandleFunc("/", s.defaultHandler)
	http.HandleFunc("/heartbeat", s.heartbeatHandler)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
