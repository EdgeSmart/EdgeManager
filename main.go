package main

import (
	"fmt"
	"os"

	"github.com/EdgeSmart/EdgeManager/service/mqtt"
	"github.com/EdgeSmart/EdgeManager/service/proxy"
)

/*
Restful API
*/
func main() {
	signal := make(chan os.Signal)

	// Start MQTT server
	go mqtt.Run()

	// Start proxy server
	// go proxy.Run()

	go proxy.RunNew()

	quitSignal := <-signal
	fmt.Println("Process quit: ", quitSignal)
}
