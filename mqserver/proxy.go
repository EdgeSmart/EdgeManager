package mqserver

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/EdgeSmart/EdgeManager/library"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

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
	for {
		// 可能的bug，如果正好第最后一次遍历取出所有数据，可能出现hang住
		reqBuf := make([]byte, 4096)
		bufLen := len(reqBuf)
		readLen, err := conn.Read(reqBuf[:])
		if err != nil {
			fmt.Println(err)
			return []byte{}, errors.New("Read data failed")
		}
		requestData = append(requestData, reqBuf[0:readLen]...)
		if readLen < bufLen {
			break
		}
	}

	return requestData, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
