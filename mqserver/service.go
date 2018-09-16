package mqserver

import (
	"fmt"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

type topicItem struct {
	topic   string
	qos     byte
	handler MQTT.MessageHandler
}

type ProxyStruct struct {
	ResponseTopic string
	Data          []byte
}

var (
	mqttClient MQTT.Client
	clientID   string
	mqttServe  chan os.Signal
	topicList  = []topicItem{
		topicItem{topic: "Manager/Monitor", qos: 0, handler: monitorHandler},
		topicItem{topic: "Manager/APP/Install", qos: 0, handler: appHandler},
		topicItem{topic: "Manager/APP/Uninstall", qos: 0, handler: appHandler},
		topicItem{topic: "Manager/APP/List", qos: 0, handler: appHandler},
		topicItem{topic: "Manager/gateway_register", qos: 0, handler: gatewayRegisterHandler},
	}
)

func startService() {
	clientID = fmt.Sprintf("manager/%s", "main")
	opts := MQTT.NewClientOptions().AddBroker("tcp://127.0.0.1:1883")
	opts.SetClientID(clientID)
	opts.SetUsername(clientID)
	opts.SetPassword("no_need")
	opts.SetDefaultPublishHandler(messageHandler)

	mqttClient = MQTT.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	for i := range topicList {
		topic := topicList[i].topic
		if token := mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
	}
}

func messageHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func monitorHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func appHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func gatewayRegisterHandler(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func Proxy(ctx *gin.Context) string {

	return ""
}
