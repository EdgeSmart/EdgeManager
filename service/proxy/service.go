package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/EdgeSmart/EdgeManager/service/mqtt"
)

var (
	serviceHost = ":2375"
)

// RunNew RunNew
func RunNew() {
	log.Println("Starting HTTP to MQTT proxy server ...")
	http.HandleFunc("/", mqtt.Proxy)
	log.Fatal(http.ListenAndServe(serviceHost, nil))
}

// Run run a proxy server
func Run() {
	service := ":8081"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		fmt.Println(err)
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go mqtt.ProxyOld(conn)
	}
}
