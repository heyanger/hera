package main

import (
	"github.com/funkytennisball/hera/datanode"
)

func main() {
	// Hardcode lol
	service := 'localhost:8080'

	dn := datanode.DataNode{service}
	dn.Init()
}
