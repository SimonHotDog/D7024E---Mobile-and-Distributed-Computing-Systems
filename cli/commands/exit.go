package commands

import (
	"d7024e/kademlia"
	"fmt"
	"os"
)

func ExitApplication(context *kademlia.Kademlia, args string) (string, error) {
	fmt.Println("Goodbye!")
	os.Exit(0)
	return "", nil
}
