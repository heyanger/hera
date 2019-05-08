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
	if len(os.Args) < 3 {
		fmt.Println("Usage go run main.go 'configfile' 'servername'")

		return
	}

	configFile := os.Args[1]
	// serverName := os.Args[2]

	dat, err := ioutil.ReadFile(configFile)

	if err != nil {
		fmt.Println("Cannot read file: " + configFile)
	}

	// m := common.ServersConfig{}
	m := Config{}
	// servers := []common.Servers{}

	yaml.Unmarshal(dat, &m)

	// for _, v := range m {
	// 	// s := common.Servers{Id: k.(string), Address: v.([]interface{})["protocolport"].(string)}
	// 	// servers = append(servers, s)
	// 	fmt.Println(v)
	// }
	// fmt.Println(m)

	r := new(protocol.SRaft)
	s := service.Service{Protocol: r}

	r.Init([]string{"1", "2", "3"}, "1")
	s.Start("8080")
}
