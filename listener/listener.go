package listener

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"
)

const topicName string = "owntracks/+/+/event"

type TransitionMessage struct {
	Wtst  int64   // Time of waypoint creation
	Lat   float32 // Latitude
	Long  float32 // Longitude
	Tst   int64   // Timestamp of transition
	Acc   uint32  // Accuracy of Lat/Long
	Tid   string  // Tracker ID
	Event string  // Enter or Leave
	Desc  string  // Description
}

//define a function for the default message handler
var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())

	res := TransitionMessage{}
	json.Unmarshal([]byte(msg.Payload()), &res)

	fmt.Printf("DESC: %s\n", res.Desc)
	time := time.Unix(res.Tst, 0)
	fmt.Printf("TIME: %s\n", time)

	fmt.Println()
}

func sampleMessage() string {
	obj := &TransitionMessage{
		Tst:   time.Now().Unix(),
		Event: "enter",
		Desc:  "Test Desc"}
	msg, _ := json.Marshal(obj)
	return string(msg)
}

func Listen(connectionString string) {
	options := mqtt.NewClientOptions().AddBroker(connectionString)
	options.SetClientID("eventr")
	options.SetDefaultPublishHandler(f)

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	if token := client.Subscribe(topicName, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	//Publish 5 messages to the topic at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 5; i++ {
		text := sampleMessage()
		token := client.Publish("owntracks/phil/iPhone/event", 0, false, text)
		token.Wait()
		time.Sleep(500 * time.Millisecond)
	}

	//unsubscribe from topic
	if token := client.Unsubscribe(topicName); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}
