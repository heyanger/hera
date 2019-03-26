package main

import (
	"fmt"

	"github.com/funkytennisball/hera/protocol"
)

func main() {
	var r protocol.Protocol = protocol.Raft{}
	var k protocol.Key = protocol.Key{K: "key"}
	var e protocol.Entity = protocol.Entity{V: "Value"}

	fmt.Println(r.Insert(k, e))
}
