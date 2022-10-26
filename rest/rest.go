package rest

import (
	"d7024e/cli/commands"
	"d7024e/kademlia"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var context kademlia.IKademlia

type Data struct {
	Data     string
	Location string
}

// Replies to successful put (store) requests with appropriate HTTP headers and json data and denies failed lookups
func putHandle(w http.ResponseWriter, r *http.Request) {
	//help from stackoverflow.com/questions/46579429/golang-cant-get-body-from-request-getbody'
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Invalid HTTP request, try /objects/{hash} for GET"))
		return
	}
	b, err := io.ReadAll(r.Body)
	if err == nil {
		str := string(b[:])
		split := strings.Split(str, "=")
		data := split[1]
		res, errPut := commands.PutObjectInStore(context, data)

		w.Header().Set("Content-Type", "application/json")
		if errPut != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusCreated)
			reply := Data{data, "/objects/" + res}
			json.NewEncoder(w).Encode(reply)
		}
	} else {
		fmt.Fprintf(w, err.Error())
	}

}

// Replies to get (lookup) requests with json data of the lookup target
func getHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Invalid HTTP request, try /objects for POST"))
		return
	}
	hash := strings.Split(r.URL.Path, "/")[2]
	str, err := commands.GetObjectByHash(context, hash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, err.Error())
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, str)
	}

}

// Homepage guide
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Example of post: /objects")
	fmt.Fprintln(w, "Example of get: /objects/{hash}")
}

// Directs webpages to corresponding handlers and starts listener
func Restful(kademlia kademlia.IKademlia) {
	context = kademlia
	http.HandleFunc("/", homePage)
	http.HandleFunc("/objects", putHandle)
	http.HandleFunc("/objects/", getHandle)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
