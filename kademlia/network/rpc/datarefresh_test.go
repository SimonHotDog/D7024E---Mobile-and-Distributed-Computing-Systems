package rpc

import (
	"d7024e/kademlia/datastore"
	"d7024e/kademlia/network"
	"d7024e/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSendRefreshDataMessage(t *testing.T) {
	startTime := time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	refreshTime := startTime.Add(time.Minute)
	dataTTL := time.Hour
	dataKey := "key"
	timeprovider := &util.FakeTimeProvider{InternalTime: startTime}
	datastore := datastore.NewDataStore(dataTTL, nil, timeprovider)
	networkB, _ := network.NewNetwork(14048, datastore)
	networkA, _ := network.CreateTestNetwork(14041)

	go networkA.Listen()
	go networkB.Listen()
	defer networkA.StopListen()
	defer networkB.StopListen()

	time.Sleep(20 * time.Millisecond)

	networkB.GetDatastore().Set(dataKey, []byte("test"))
	timeprovider.InternalTime = refreshTime

	SendRefreshDataMessage(networkA, networkB.GetMe(), "hash")
	timeprovider.InternalTime = startTime.Add(dataTTL)

	// Data should be in datastore.
	// If refresh failed, data should not be in datastore
	_, actualExists := datastore.Get(dataKey)

	assert.True(t, actualExists)
}
