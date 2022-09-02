package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
)

func main() {

	context := kademlia.Kademlia{}

	cli.PrintHello()
	cli.Open(&context)
}
