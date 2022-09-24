package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

func init() {
	log.SetFlags(log.Ltime | log.Lshortfile)
}

func main() {
	port, bootstrapAddress, verbose := retriveProgramParameters()
	if !*verbose {
		log.SetOutput(io.Discard)
	}

	bootstrap := kademlia.NewContact(nil, *bootstrapAddress)
	myAddress := fmt.Sprintf("%s:%d", kademlia.GetOutboundIP(), *port)
	me := kademlia.NewContact(kademlia.NewRandomKademliaID(), myAddress)
	context := kademlia.Kademlia{Me: &me, Routing: kademlia.NewRoutingTable(me)}
	network := kademlia.CreateNewNetwork(&context, *port)

	context.Network = &network

	go network.Listen() // TODO: Notify it is actually listening
	time.Sleep(1 * time.Second)
	go context.JoinNetwork(&bootstrap)
	cli.Open(&context)
}

func retriveProgramParameters() (port *int, bootstrapNode *string, verbose *bool) {
	env_port, _ := strconv.Atoi(os.Getenv("KADEMLIA_PORT"))
	env_verbose, _ := strconv.ParseBool(os.Getenv("KADEMLIA_VERBOSE"))
	env_bootstrapNode := os.Getenv("KADEMLIA_BOOTSTRAP_NODE")

	port = flag.Int("p", env_port, "Portnumber")
	verbose = flag.Bool("v", env_verbose, "Indicates if a log should be created")
	bootstrapNode = flag.String("b", env_bootstrapNode, "Adress of bootstrap node")

	flag.Parse()

	return port, bootstrapNode, verbose
}
