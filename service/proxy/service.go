package proxy

import (
	"log"
	"net/http"

	"github.com/EdgeSmart/EdgeManager/service/mqtt"
)

var (
	serviceHost = ":2375"
)

// Run run proxy
func Run() {
	log.Println("Starting HTTP to MQTT proxy server ...")
	http.HandleFunc("/", mqtt.Proxy)
	log.Fatal(http.ListenAndServe(serviceHost, nil))
}
