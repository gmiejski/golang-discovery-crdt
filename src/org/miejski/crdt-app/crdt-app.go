package main

import (
	"org/miejski/discovery"
	"fmt"
	"net/http"
)

func main() {


	discovery_client := discovery.NewDiscoveryClient()

	discovery_client.CurrentActiveNodes()
	fmt.Println(discovery_client.CurrentActiveNodes())

	http.HandleFunc("/status", handleRequest)
	http.ListenAndServe(":8080", nil)
}

func handleRequest(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":

	}
}
