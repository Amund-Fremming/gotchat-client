package main

import (
	"client/cmd"
	"client/program"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/amund-fremming/gotchat-common/config"
)

func main() {
	config, err := config.Load()
	if err != nil {
		fmt.Print("\r[CLIENT] Failed due to missing enviroment variables\n")
		return
	}

	serverUrlBase := fmt.Sprintf("%s://%s:%s", config.Scheme, config.URL, config.Port)
	slog.SetLogLoggerLevel(config.LogLevel)

	response, err := http.Get(serverUrlBase + "/health")
	if err != nil || response.Status != "200 OK" {
		fmt.Println("[ERROR] The sever is currently unavailable")
		fmt.Println("[CLIENT] Shutting down..")
		return
	}
	fmt.Println("[SERVER] Healthy")

	err = cmd.InitReadline()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	program.ConnectToServer(&config)
	cmd.DisplayWelcomeMessage()

	go program.CommandReader()
	go program.ServerReader()
	go program.CommandDispatcher()

	select {}
}
