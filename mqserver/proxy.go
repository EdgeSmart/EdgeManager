package mqserver

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/EdgeSmart/EdgeManager/library"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/vmihailenco/msgpack"
)

type Proxy struct {
	ResponseTopic string
	Method        string
	URL           string
	Proto         string
	Host          string
	Header        http.Header
	Body          []byte
}

// HandleClientV2 HandleClientV2
func HandleClientV2(ctx *gin.Context) {
	responseChan := make(chan []byte)
	randTopicStr := library.GetRandomString(32)
	topic := "Manager/docker_proxy/response_" + randTopicStr
	proxyHandler := func(client MQTT.Client, msg MQTT.Message) {
		responseChan <- msg.Payload()
		mqttClient.Unsubscribe(topic)
	}
	if token := mqttClient.Subscribe(topic, 0, proxyHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	body, _ := ctx.GetRawData()

	proxyData := Proxy{
		ResponseTopic: topic,
		Method:        ctx.Request.Method,
		URL:           ctx.Request.URL.String(),
		Proto:         ctx.Request.Proto,
		Host:          ctx.Request.Host,
		Header:        ctx.Request.Header,
		Body:          body,
	}

	proxyDataBytes, _ := msgpack.Marshal(&proxyData)

	// publish register message
	token := mqttClient.Publish("cluster_test/DockerAPI_test", 0, false, proxyDataBytes)
	token.Wait()
}

func HandleClient(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	defer conn.Close()
	requestData, err := readData(conn)
	if err != nil {
		return
	}

	responseChan := make(chan []byte)
	randTopicStr := library.GetRandomString(32)
	topic := "Manager/docker_proxy/response_" + randTopicStr
	proxyHandler := func(client MQTT.Client, msg MQTT.Message) {
		responseChan <- msg.Payload()
		mqttClient.Unsubscribe(topic)
	}
	if token := mqttClient.Subscribe(topic, 0, proxyHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	proxyData := ProxyStruct{
		ResponseTopic: topic,
		Data:          requestData,
	}

	var GolngGob bytes.Buffer
	enc := gob.NewEncoder(&GolngGob)

	err = enc.Encode(&proxyData)
	if err != nil {
		log.Fatal("encode error:", err)
	}
	data := GolngGob.Bytes()

	// publish register message
	token := mqttClient.Publish("cluster_test/DockerAPI", 0, false, data)
	token.Wait()

	response := <-responseChan
	conn.Write(response)
	conn.Close()
}

func readData(conn net.Conn) ([]byte, error) {
	requestData := []byte{}
	requestData, err := ioutil.ReadAll(conn)
	fmt.Println(requestData)
	fmt.Println(string(requestData))
	return requestData, err
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
