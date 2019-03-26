package main

import (
	"github.com/funkytennisball/hera/protocol"
	"github.com/funkytennisball/hera/service"
)

func main() {
	r := protocol.Raft{}
	s := service.Service{Protocol: r}
	s.Start()
}
