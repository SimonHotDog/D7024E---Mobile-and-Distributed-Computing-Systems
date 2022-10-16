package rpc

import (
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network"
	"d7024e/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendForgetDataMessage(t *testing.T) {
	dataTTL := time.Hour
	dataKey := "key"
	timeprovider := new(util.TimeProvider)
	datastore := datastore.NewDataStore(dataTTL, nil, timeprovider)
	networkB, _ := network.NewNetwork(14048, datastore)
	networkA, _ := network.CreateTestNetwork(14041)

	go networkA.Listen()
	go networkB.Listen()
	defer networkA.StopListen()
	defer networkB.StopListen()

	time.Sleep(20 * time.Millisecond)

	networkB.GetDatastore().Set(dataKey, []byte("test"))

	SendForgetDataMessage(networkA, networkB.GetMe(), dataKey)

	time.Sleep(20 * time.Millisecond)

	_, actualExists := datastore.Get(dataKey)

	assert.False(t, actualExists)
}
