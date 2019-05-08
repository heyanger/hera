package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/funkytennisball/hera/protocol"
	"github.com/funkytennisball/hera/service"
	"gopkg.in/yaml.v2"
)

type Server struct {
	Clientport   string
	Protocolport string
}

type Config struct {
	Servers map[string]string
}

func main() {
	configFile := "./config.yaml"

	if len(os.Args) >= 2 {
		configFile = os.Args[1]
	}

	dat, err := ioutil.ReadFile(configFile)

	if err != nil {
		fmt.Println("Cannot read file: " + configFile)
	}

	m := Config{}

	yaml.Unmarshal(dat, &m)

	r := new(protocol.Raft)
	s := service.Service{Protocol: r}

	r.Init([]string{"1", "2", "3"}, "1")
	s.Start("8080")
}
