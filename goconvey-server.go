package main

import (
	"flag"

	_ "github.com/jacobsa/oglematchers"
	_ "github.com/jacobsa/ogletest"
)

func init() {
	flag.IntVar(&port, "port", 8080, "The port at which to serve http.")
	flag.StringVar(&host, "host", "127.0.0.1", "The host at which to serve http.")
}

func main() {
	flag.Parse()

	// initialize file system
	// initialize shell

	// initialize a watcher
	// adjust watcher

	// initialize parser
	// initialize executor
	// initialize server

	// initialize scanner
	// start scanning

	// serve http
}

var (
	port int
	host string
)
