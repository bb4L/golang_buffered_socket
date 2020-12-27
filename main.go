package main

import (
	"golangbufferedsocket/logging"
	"golangbufferedsocket/restserver"
	"io/ioutil"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type config struct {
	Connections []restserver.Config `yaml:"connections"`
}

func main() {

	argsWithoutProg := os.Args[1:]
	logging.GetInfoLogger().Printf("Starting server with %s", argsWithoutProg)

	if len(argsWithoutProg) != 1 {
		logging.GetFatalLogger().Fatalln("wrong number of arguments, exactly 1 is needed")
	}
	filename := os.Args[1]
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var servers config

	source := []byte(data)
	err = yaml.Unmarshal(source, &servers)

	for _, server := range servers.Connections {
		logging.GetInfoLogger().Println("Starting server loop")
		logging.GetInfoLogger().Println(server)
		// TODO: currently multiple socket connections aren't working
		runServerConfig(server)
	}
}

func runServerConfig(config restserver.Config) {
	logging.GetInfoLogger().Println("Starting server..")
	logging.GetInfoLogger().Println(config)
	if config.ServerType == "" {
		config.ServerType = "unix"
	}

	if config.ServerAddress == "" {
		logging.GetErrorLogger().Printf("No server address provided")
		return
	}

	l, err := net.Listen(config.ServerType, config.ServerAddress)
	if err != nil {
		logging.GetErrorLogger().Fatal("listen error:", err)
	}

	for {
		fd, err := l.Accept()
		if err != nil {
			logging.GetErrorLogger().Fatal("accept error:", err)
		}

		logging.GetInfoLogger().Println("Start listening to connection")
		go restserver.Server(fd, config)
	}
}
