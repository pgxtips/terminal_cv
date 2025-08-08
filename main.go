package main

import (
	"flag"
	"terminal_cv/cmd"
)

func main(){
	var host = flag.String("host", "127.0.0.1", "Host address for SSH server to listen")
	var port = flag.Int("port", 1337, "Port for SSH server to listen")

	flag.Parse()

	terminal_cv.StartServer(*host, *port)
}
