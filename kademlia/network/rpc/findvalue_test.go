package rpc

import (
	"d7024e/kademlia/network"
	"d7024e/util"
	"strings"
	"testing"
	"time"
)

func TestLookupMessage(t *testing.T) {
	testname := "Send lookup command"
	t.Run(testname, func(t *testing.T) {
		expected := "My Message"
		valueChannel := make(chan string, 1)
		networkA, _ := network.CreateTestNetwork(14041)
		networkB, _ := network.CreateTestNetwork(14048)
		messageBytes := []byte(expected)
		messageHash := util.Hash([]byte(messageBytes))
		networkB.Datastore.Set(messageHash, messageBytes)

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)

		go SendLookupMessage(networkA, networkB.Me, messageHash, valueChannel)
		actual := <-valueChannel

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestLookupMessageTimeout(t *testing.T) {
	testname := "Send lookup command and wait for timeout"
	t.Run(testname, func(t *testing.T) {
		expected := network.NETWORK_REQUEST_TIMEOUT_STRING
		valueChannel := make(chan string, 1)
		networkA, _ := network.CreateTestNetwork(14041)
		networkB, _ := network.CreateTestNetwork(14048)
		messageBytes := []byte("Test")
		messageHash := util.Hash([]byte(messageBytes))
		networkB.Datastore.Set(messageHash, messageBytes)

		go networkA.Listen()
		defer networkA.StopListen()
		time.Sleep(20 * time.Millisecond)

		go SendLookupMessage(networkA, networkB.Me, messageHash, valueChannel)
		actual := <-valueChannel

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestStoreBeforeLookupMessage(t *testing.T) {
	testname := "Send store message before lookup"
	t.Run(testname, func(t *testing.T) {
		expected := "My Message"
		valueChannel := make(chan string, 1)
		networkA, _ := network.CreateTestNetwork(14041)
		networkB, _ := network.CreateTestNetwork(14048)
		messageBytes := []byte(expected)
		messageHash := util.Hash([]byte(messageBytes))

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)

		storeOk := SendStoreMessage(networkA, networkB.Me, messageHash, messageBytes)
		if !storeOk {
			t.Errorf("Expected store to succeed")
			return
		}
		go SendLookupMessage(networkA, networkB.Me, messageHash, valueChannel)
		actual := <-valueChannel

		a, oka := networkA.Datastore.Get(messageHash)
		b, okb := networkB.Datastore.Get(messageHash)
		_ = a
		_ = oka
		_ = b
		_ = okb

		if strings.Compare(actual, expected) != 0 {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}
