package main

import (
	"fmt"
	"os"

	"github.com/funkytennisball/hera/protocol"
	"github.com/funkytennisball/hera/service"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage go run main.go 'serverid' 'port'")

		return
	}

	id := os.Args[1]
	port := os.Args[2]

	r := new(protocol.Raft)
	s := service.Service{Protocol: r}

	r.Init([]string{"1", "2", "3"}, id)
	s.Start(port)
}
