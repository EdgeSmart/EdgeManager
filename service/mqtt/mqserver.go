package mqtt

import (
	"time"

	"github.com/surgemq/surgemq/service"
)

var server = &service.Server{
	KeepAlive:        300,         // seconds
	ConnectTimeout:   2,           // seconds
	SessionsProvider: "mem",       // keeps sessions in memory
	Authenticator:    "edge_auth", // always succeed
	TopicsProvider:   "mem",       // keeps topic subscriptions in memory
}

// Run run mqtt server
func Run() {
	// Listen and serve connections at localhost:1883
	go server.ListenAndServe("tcp://:1883")
	time.Sleep(200000)
	startService()
}
