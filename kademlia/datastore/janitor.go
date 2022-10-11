package datastore

// Inspiration: https://github.com/patrickmn/go-cache/

import "time"

type Janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *Janitor) Run(store *DataStore) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			store.RemoveExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(store *DataStore) {
	store.janitor.stop <- true
}

// Create a new janitor with the given interval.
func runJanitor(store *DataStore, ci time.Duration) {
	j := &Janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	store.janitor = j
	go j.Run(store)
}
