package datanode

import (
	"encoding/json"
	"net/http"
)

// DataNode defines the lifecycle of a data node
type DataNode struct {
	service string
}

// Init starts the Data Node
func (d *DataNode) Init() {
	values := map[string]string{}

}

func (d *DataNode) request(values) {
	jsonValue, _ := json.Marshal(values)
	resp, err := http.Post(d.service, "application/json", jsonValue)
}
