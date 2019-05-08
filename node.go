package main

import (
	"fmt"
	"os"

	"github.com/funkytennisball/hera/datanode"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("usage: go run node.go PORT RAFTPORT")
		return
	}

	port := os.Args[1]
	raftport := os.Args[2]

	// Hardcode lol
	service := "http://localhost:8080"

	dn := datanode.DataNode{Service: service}
	dn.Init(port, raftport)
}
