package commands

import (
	"d7024e/kademlia"
	"d7024e/kademlia/network/routing"
	"d7024e/kademlia/network/rpc"
)

func Debug_sendPing(context kademlia.IKademlia, args string) (string, error) {
	contact := routing.Contact{Address: args, ID: nil}
	alive := make(chan bool)
	go rpc.SendPingMessage(context.GetNetwork(), &contact, alive)

	if <-alive {
		return "Node is alive", nil
	} else {
		return "Node is dead", nil
	}
}
