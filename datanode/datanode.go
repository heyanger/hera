package datanode

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/funkytennisball/hera/protocol"

	"github.com/funkytennisball/hera/common"
)

// DataNode defines the lifecycle of a data node
type DataNode struct {
	Service  string
	Protocol common.Protocol

	id     string
	source string
}

// Init starts the Data Node
func (d *DataNode) Init(port string, raftport string) {
	d.Protocol = new(protocol.Raft)
	d.source = "http://127.0.0.1:" + port

	notify := make(map[string]string)
	notify["source"] = d.source

	resp := d.request("/heartbeat", notify)

	body, _ := ioutil.ReadAll(resp.Body)

	m := make(map[string]string)
	json.Unmarshal(body, &m)

	// Print Dict
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Print(string(b))

	isGenesis, _ := strconv.ParseBool(m["genesis"])
	id := m["id"]
	d.id = id

	if isGenesis {
		d.Protocol.Init(id, raftport, true)
	} else {
		n := make(map[string]string)
		n["id"] = id
		n["source"] = "127.0.0.1:" + raftport

		d.Protocol.Init(id, raftport, false)
		d.extrequest(m["leadersource"], "/join", n)
	}

	d.heartbeat()

	http.HandleFunc("/get", d.get)
	http.HandleFunc("/put", d.put)
	http.HandleFunc("/delete", d.delete)
	http.HandleFunc("/join", d.join)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func (d *DataNode) heartbeat() {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				notify := make(map[string]string)
				notify["key"] = d.id
				notify["source"] = d.source
				notify["location"] = d.Protocol.Location()
				notify["leader"] = strconv.FormatBool(strings.Compare(d.Protocol.State(), "Leader") == 0)

				d.request("/heartbeat", notify)

				// resp := d.request("/heartbeat", notify)
				// // Printing
				// body, _ := ioutil.ReadAll(resp.Body)
				// m := make(map[string]string)
				// json.Unmarshal(body, &m)
				// b, _ := json.MarshalIndent(m, "", "  ")
				// fmt.Print(string(b))
			}
		}
	}()
}

func (d *DataNode) get(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Write([]byte(d.Protocol.Get(common.Key{K: m["key"]}).V))
}

func (d *DataNode) put(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := d.Protocol.Insert(common.Key{K: m["key"]}, common.Entity{V: m["value"]}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (d *DataNode) delete(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// TODO: error handling
	if err := d.Protocol.Remove(common.Key{K: m["key"]}); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (d *DataNode) join(w http.ResponseWriter, r *http.Request) {
	m := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d.Protocol.Join(m["id"], m["source"])

	w.WriteHeader(http.StatusOK)
}

func (d *DataNode) request(path string, values map[string]string) *http.Response {
	jsonValue, _ := json.Marshal(values)
	resp, _ := http.Post(d.Service+path, "application/json", bytes.NewBuffer(jsonValue))

	return resp
}

func (d *DataNode) extrequest(source string, path string, values map[string]string) *http.Response {
	jsonValue, _ := json.Marshal(values)
	resp, _ := http.Post(source+path, "application/json", bytes.NewBuffer(jsonValue))

	return resp
}
