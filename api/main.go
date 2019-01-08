package main

import (
	api "github.com/mlh2018/api/api_server"
	"os"
	"strconv"
)

func main() {
	var appID = os.Getenv("AppID")
	var endpointKey = os.Getenv("EndpointKey")
	var region = os.Getenv("Region")
	var secret = os.Getenv("Secret")
	var opcKey = os.Getenv("OcpKey")
	var port = os.Getenv("Port")
	portNumb, err := strconv.Atoi(port)
	if portNumb == 0 || err != nil {
		panic("invalid port number")
	}
	client := api.NewClient(appID, endpointKey, region, opcKey)
	s := api.NewAPIServer(portNumb, secret, client)
	s.Run()
}
