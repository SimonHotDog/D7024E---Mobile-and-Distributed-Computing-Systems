package main

import (
	"d7024e/cli"
	"d7024e/kademlia"
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network"
	"d7024e/kademlia/network/routing"
	"d7024e/rest"
	"d7024e/util"
	"flag"
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

	timeprovider := &util.TimeProvider{}
	datastore := datastore.NewDataStore(time.Hour, nil, timeprovider)
	bootstrap := routing.NewContact(nil, *bootstrapAddress)
	network, me := network.NewNetwork(*port, datastore)
	context := kademlia.NewKademlia(me, network, datastore)
	cli := cli.NewCli(os.Stdout, os.Stdin, context)

	go network.Listen() // TODO: Notify it is actually listening
	time.Sleep(1 * time.Second)
	go context.JoinNetwork(&bootstrap, 60)
	go rest.Restful(context)
	cli.Open(true)
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
