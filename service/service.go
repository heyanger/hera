package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/funkytennisball/hera/common"
)

// Service provides an HTTP service
type Service struct {
	ranges common.RangeMap
	nodes  common.NodeMap

	leader common.NodeKey
}

func (s *Service) postHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notify := map[string]string{}
	notify["key"] = m["key"]

	resp := s.request(s.getLeaderSource(), "/get", notify)

	body, _ := ioutil.ReadAll(resp.Body)

	w.Write([]byte(body))
}

func (s *Service) putHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	notify := map[string]string{}
	notify["key"] = m["key"]
	notify["value"] = m["value"]

	s.request(s.getLeaderSource(), "/put", notify)

	w.WriteHeader(http.StatusOK)
}

func (s *Service) getLeaderSource() string {
	return s.nodes[s.leader].Source
}

func (s *Service) deleteHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) newNode(m map[string]string) map[string]string {
	res := make(map[string]string)

	isGenesis := len(s.nodes) == 0

	id := rand.Uint64()
	key := common.NodeKey(id)
	v := common.Node{
		Source:    m["source"],
		Heartbeat: uint64(time.Now().UnixNano() / 1000),
	}
	s.nodes[key] = v

	res["id"] = strconv.FormatUint(id, 10)
	res["genesis"] = strconv.FormatBool(isGenesis)
	res["source"] = m["source"]

	if s.leader != 0 {
		res["leaderid"] = strconv.FormatUint(uint64(s.leader), 10)
		res["leaderaddr"] = s.nodes[s.leader].Location
	}

	return res
}

func (s *Service) request(leader string, path string, values map[string]string) *http.Response {
	jsonValue, _ := json.Marshal(values)
	resp, _ := http.Post(leader+path, "application/json", bytes.NewBuffer(jsonValue))

	return resp
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
		jData, _ := json.Marshal(s.newNode(m))
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return
	}

	key := common.NodeKey(i)

	if _, ok := s.nodes[key]; ok {
		// Node exists, update heartbeat
		v := s.nodes[key]
		v.Heartbeat = uint64(time.Now().UnixNano() / 1000)
		v.Source = m["source"]
		v.Location = m["location"]
		v.Leader, _ = strconv.ParseBool(m["leader"])
		s.nodes[key] = v

		res := make(map[string]string)

		if v.Leader {
			s.leader = key

			res["leaderid"] = strconv.FormatUint(uint64(s.leader), 10)
			res["leaderaddr"] = s.nodes[s.leader].Location
		}

		jData, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return
	} else {
		// Node does not exist, add new entry
		jData, _ := json.Marshal(s.newNode(m))
		w.Header().Set("Content-Type", "application/json")
		w.Write(jData)
		return
	}
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
	s.nodes = make(common.NodeMap)
	s.ranges = make(common.RangeMap)

	http.HandleFunc("/", s.defaultHandler)
	http.HandleFunc("/heartbeat", s.heartbeatHandler)

	fmt.Println("Main Server Started on port: " + port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
