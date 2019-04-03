package main

import (
	"github.com/funkytennisball/hera/protocol"
	"github.com/funkytennisball/hera/service"
)

func main() {
	r := new(protocol.Raft)
	s := service.Service{Protocol: r}

	r.Init(true, "1")
	s.Start()
}
