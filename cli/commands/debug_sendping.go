package commands

import "d7024e/kademlia"

func Debug_sendPing(context *kademlia.Kademlia, args string) (string, error) {
	contact := kademlia.Contact{Address: args, ID: nil}
	alive := make(chan bool)
	go context.Network.SendPingMessage(&contact, alive)

	if <-alive {
		return "Node is alive", nil
	} else {
		return "Node is dead", nil
	}
}
