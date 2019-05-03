package mqtt

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EdgeSmart/EdgeManager/library"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type ProxyData struct {
	ResponseTopic string
	Method        string
	URL           string
	Proto         string
	Host          string
	Header        http.Header
	Body          []byte
}

const (
	lineBreak = "\r\n"
	baseHost  = "dapi.hopeness.net"
	bufLen    = 1024
)

// Proxy HTTP to MQTT proxy
func Proxy(w http.ResponseWriter, r *http.Request) {
	var (
		httpCode = 502
		resData  = []byte("{\"message\": \"proxy server inner error\"}")
		header   = w.Header()
	)
	// defer the response
	defer func() {
		header := w.Header()
		header["Content-Type"] = []string{"application/json"}
		w.WriteHeader(httpCode)
		w.Write(resData)
	}()
	// get request raw data
	buf := bytes.NewBuffer([]byte{})
	err := r.WriteProxy(buf)
	if err != nil {
		log.Println("Proxy data bufer error")
		return
	}
	requestData := buf.Bytes()
	// check the cluster name
	// todo: select db to verify exist and check state if online
	clusterName := r.Host[:strings.Index(r.Host, ".")]
	if len(clusterName) == 0 {
		log.Println("Cluster name error")
		return
	}
	responseChan := make(chan []byte)
	// subscribe
	randTopicStr := library.GetRandomString(32)
	topic := "Manager/docker_proxy/response_" + randTopicStr
	proxyHandler := func(client MQTT.Client, msg MQTT.Message) {
		responseChan <- msg.Payload()
		mqttClient.Unsubscribe(topic)
	}
	if token := mqttClient.Subscribe(topic, 0, proxyHandler); token.Wait() && token.Error() != nil {
		log.Println("error")
		return
	}
	proxyData := ProxyStruct{
		ResponseTopic: topic,
		Data:          requestData,
	}

	var GolngGob bytes.Buffer
	enc := gob.NewEncoder(&GolngGob)

	err = enc.Encode(&proxyData)
	if err != nil {
		log.Println("encode error:", err)
		return
	}
	data := GolngGob.Bytes()
	// publish register message
	token := mqttClient.Publish(clusterName+"/DockerAPI", 0, false, data)
	token.Wait()
	// get response data
	response := <-responseChan
	resStr := string(response)
	resHeaderStr := strings.Trim(resStr[:strings.Index(resStr, lineBreak+lineBreak)], " "+lineBreak)
	resBodyStr := strings.Trim(resStr[strings.Index(resStr, lineBreak+lineBreak):], " "+lineBreak)
	resHeaderMap := strings.Split(resHeaderStr, lineBreak)
	headerExcept := map[string]bool{
		"Content-Length": true,
	}
	// handle header
	for i, keyVal := range resHeaderMap {
		if i == 0 {
			httpScheme := strings.Split(keyVal, " ")
			httpCode, err = strconv.Atoi(httpScheme[1])
			if err != nil {
				log.Println("error")
				return
			}
			continue
		}
		sepPos := strings.Index(keyVal, ":")
		key := strings.Trim(keyVal[:sepPos], " ")
		value := strings.Trim(keyVal[sepPos+1:], " ")
		if _, exists := headerExcept[key]; exists {
			continue
		}
		header[key] = []string{value}
	}
	resData = []byte(resBodyStr)
	log.Printf("http2mqtt proxy, cluster: %s, topic: %s, uri: %s", clusterName, topic, r.RequestURI)
}

// ProxyOld HTTP to MQTT proxy
func ProxyOld(conn net.Conn) {
	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	defer conn.Close()
	requestData, err := readData(conn)
	if err != nil {
		return
	}
	// check the cluster name
	// todo: select db to verify exist and check state if online
	reqLines := strings.Split(string(requestData), "\r")
	clusterName := ""
	for _, line := range reqLines {
		if endPos := strings.Index(line, baseHost); endPos > -1 {
			clusterName = line[7 : endPos-1]
		}
	}
	if len(clusterName) == 0 {
		log.Fatal("error")

	}
	responseChan := make(chan []byte)
	// subscribe
	randTopicStr := library.GetRandomString(32)
	topic := "Manager/docker_proxy/response_" + randTopicStr
	proxyHandler := func(client MQTT.Client, msg MQTT.Message) {
		responseChan <- msg.Payload()
		mqttClient.Unsubscribe(topic)
	}
	if token := mqttClient.Subscribe(topic, 0, proxyHandler); token.Wait() && token.Error() != nil {
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
	token := mqttClient.Publish(clusterName+"/DockerAPI", 0, false, data)
	token.Wait()
	// get response data
	response := <-responseChan
	conn.Write(response)
	conn.Close()
}

func readData(conn net.Conn) ([]byte, error) {
	var (
		buffer      bytes.Buffer
		err         error
		requestData []byte
	)
	for {
		buf := make([]byte, bufLen)
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		buffer.Write(buf[:n])
		if n < bufLen {
			break
		}
	}
	requestData = buffer.Bytes()
	return requestData, err
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
