package rpc

import (
	"d7024e/kademlia/network"
	"testing"
	"time"
)

func TestPingMessage(t *testing.T) {
	testname := "Ping myself"
	t.Run(testname, func(t *testing.T) {
		expected := true
		alive := make(chan bool, 1)
		networkA, _ := network.CreateTestNetwork(14041)
		networkB, _ := network.CreateTestNetwork(14048)

		go networkA.Listen()
		go networkB.Listen()
		defer networkA.StopListen()
		defer networkB.StopListen()
		time.Sleep(20 * time.Millisecond)
		SendPingMessage(networkA, networkB.GetMe(), alive)
		actual := <-alive

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}

func TestPingMessageFailure(t *testing.T) {
	testname := "Fail to ping myself"
	t.Run(testname, func(t *testing.T) {
		expected := false
		alive := make(chan bool, 1)
		network, me := network.CreateTestNetwork(14041)

		SendPingMessage(network, me, alive)
		actual := <-alive

		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
		}
	})
}
