package main

import "net/http"

var servers []Server = make([]Server, 0)

func main() {

	// endpoint to register server
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		servers = append(servers, Server{Url: r.RemoteAddr, Weight: 1, State: 1})
	})

	// research security

	// create http server
	http.ListenAndServe(":8080", nil)

	// http using time.Ticker client to send pings
}

type Server struct {
	Url    string
	Weight int
	State  int
}
