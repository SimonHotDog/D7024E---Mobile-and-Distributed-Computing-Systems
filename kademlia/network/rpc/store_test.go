package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/util"
	"strings"
	"testing"
	"time"
)

func TestStoreMessage(t *testing.T) {
	testname := "Send store command"
	t.Run(testname, func(t *testing.T) {
		expected := "My Message"
		networkA, _ := network.CreateTestNetwork(14041)
		networkB, _ := network.CreateTestNetwork(14048)
		messageBytes := []byte(expected)
		messageHash := util.Hash([]byte(messageBytes))

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)

		SendStoreMessage(networkA, networkB.Me, messageHash, messageBytes)
		storedMessage, _ := networkB.Datastore.Get(messageHash)
		actual := string(storedMessage)

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestStoreMessageTimeout(t *testing.T) {
	testname := "Send store command and wait for timeout"
	t.Run(testname, func(t *testing.T) {
		expected := false
		networkA, _ := network.CreateTestNetwork(14041)
		networkB, _ := network.CreateTestNetwork(14048)
		messageBytes := []byte("My Message")
		messageHash := util.Hash([]byte(messageBytes))

		go networkA.Listen()
		defer networkA.StopListen()
		time.Sleep(20 * time.Millisecond)

		actual := SendStoreMessage(networkA, networkB.Me, messageHash, messageBytes)

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}
