package rpc

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type RPCPayload struct {
	Name string
	Data string
}

func RPCListen(rpcPort string) error {
	log.Println("Starting RPC server on port", rpcPort)
	listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", rpcPort))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		go rpc.ServeConn(rpcConn)
	}
}
