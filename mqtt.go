package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	mqttOpts   *mqtt.ClientOptions
	mqttClient mqtt.Client
)

const (
	topicRequest = "/weather/request"
	topicPush    = "/weather/update"
	ClientID     = "miot-weather"

	mqttPort = 1883
)

func mqttInit() error {
	mqttAddr, err := url.Parse(fmt.Sprintf("tcp://%s:%d", os.Getenv("MQTT_HOSTNAME"), mqttPort))
	if err != nil {
		panic(err)
	}

	mqttOpts = mqtt.NewClientOptions()
	mqttOpts.Servers = []*url.URL{mqttAddr}
	mqttOpts.ClientID = ClientID
	mqttOpts.Username = os.Getenv("MQTT_USERNAME")
	mqttOpts.Password = os.Getenv("MQTT_PASSWORD")
	mqttOpts.OnConnect = func(c mqtt.Client) {
		mqttClient.Subscribe(topicRequest, 0, func(_ mqtt.Client, msg mqtt.Message) {
			fmt.Fprintf(os.Stdout, "mqtt: recv. %s/%d\n", msg.Topic(), msg.MessageID())
			go func() { mqttPublish() }()
		})
	}

	mqttClient = mqtt.NewClient(mqttOpts)

	fmt.Fprintf(os.Stdout, "mqtt: connecting to %s\n", mqttOpts.Servers[0])
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		err := token.Error()
		fmt.Fprintf(os.Stderr, "mqtt: failed to connect. %s\n", err)
		return err
	}

	fmt.Fprintf(os.Stdout, "mqtt: connected\n")
	mqttPublish()

	return nil
}

func mqttPublish() {
	if mqttClient == nil {
		return
	}

	w := weatherGet()

	bytes, err := json.Marshal(w)
	if err == nil {
		token := mqttClient.Publish(topicPush, 0, false, bytes)
		token.Wait()
		err = token.Error()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "mqtt: failed to publish. %s\n", err)
	}

	fmt.Fprintf(os.Stdout, "mqtt: published msg to %s\n", topicPush)
}
