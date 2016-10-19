package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type addressVar []string

var (
	interval        int
	addressesToPing addressVar
	listenPort      int
)

func (a *addressVar) String() string {
	return fmt.Sprint(*a)
}

func (a *addressVar) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func pingAddress(address string) {
	response, err := http.Get("http://" + address + "/ping")

	if err != nil {
		log.Printf("Couldn't ping %s: %s", address, err)
	} else {
		defer response.Body.Close()
		var pingAnswer string
		json.NewDecoder(response.Body).Decode(&pingAnswer)
		log.Printf("Pinged %s, got answer: %s", address, pingAnswer)
	}
}

func startHTTPServer() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode("pong")
	})
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", listenPort), nil))
}

func init() {
	flag.IntVar(&interval, "interval", 10, "The interval between two pings in seconds")
	flag.Var(&addressesToPing, "address", "The address which should be pinged. Format <IP>:<port>")
	flag.IntVar(&listenPort, "listenPort", 8080, "The port to listen on")
}

func main() {
	flag.Parse()
	log.Printf("found some flags: interval=%v, address=%s, listenPort=%v", interval, addressesToPing, listenPort)

	go startHTTPServer()

	c := time.Tick(time.Duration(interval) * time.Second)
	for range c {
		for _, address := range addressesToPing {
			pingAddress(address)
		}
	}
}
