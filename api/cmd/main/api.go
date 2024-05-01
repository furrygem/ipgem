package main

import (
	"flag"

	"github.com/furrygem/ipgem/api/internal/logger"
	"github.com/furrygem/ipgem/api/internal/server"
)

var listenAddr string

func init() {
	flag.StringVar(&listenAddr, "listen-addr", "127.0.0.1:9000", "Specify listen address")
}

func main() {
	flag.Parse()
	s := server.NewServer(listenAddr)
	logger := logger.GetLogger()
	logger.Infof("Starting server on laddr: `%s`", listenAddr)
	logger.Fatal(s.Start())
}
