package mqtt

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/EdgeSmart/EdgeFairy/library/utils"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// ProxyData ProxyData
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
		httpCode     = 502
		resData      = []byte("{\"message\": \"proxy server inner error\"}")
		header       = w.Header()
		responseChan = make(chan []byte)
	)
	// defer the response
	defer func() {
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
	clusterName, err := getClusterName(r.Host)
	if err != nil {
		log.Println(err)
		return
	}
	// subscribe
	randTopicStr := utils.GetRandomString(32)
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
	// fix docker api response data with "f82" or "f83" ahead and "0" behind
	resBodyStr = strings.Trim(resBodyStr, "f0238")
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

func getClusterName(host string) (string, error) {
	// todo: select db to verify exist and check state if online
	clusterName := host[:strings.Index(host, ".")]
	if len(clusterName) == 0 {
		return "", errors.New("Cluster name error")
	}
	return clusterName, nil
}
