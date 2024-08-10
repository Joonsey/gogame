package main

import (
	"flag"
)

func main() {
	is_server := flag.String("server", "y", "run server")
	server_ip := flag.String("ip", "172.20.10.2", "ip")

	flag.Parse()

	if *is_server == "y" {
		RunServer()
	} else {
		RunClient(*server_ip)
	}
}
