package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

func (app *Config) rpcListen() error {
	log.Println("start listen rpc on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		go rpc.ServeConn(conn)
	}
}
