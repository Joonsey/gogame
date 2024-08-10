package main

import (
	"flag"
)

func main() {
	is_server := flag.String("server", "y", "run server")

	flag.Parse()

	if *is_server == "y" {
		RunServer()
	} else {
		RunClient()
	}
}
