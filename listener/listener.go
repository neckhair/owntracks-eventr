package listener

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Listener struct {
	TopicName      string
	Url            string
	OutputFilename string
}

func NewListener(config *Configuration) *Listener {
	return &Listener{
		TopicName:      "owntracks/+/+/event",
		Url:            config.Url,
		OutputFilename: config.Filename}
}

func (l *Listener) MessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		tm, err := NewTransitionMessage(msg.Payload())
		if err != nil {
			// TODO replace with logfile
			panic(err)
		}

		f, err := os.OpenFile(l.OutputFilename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			// TODO replace with logfile
			panic(err)
		}

		defer f.Close()

		lineToWrite := fmt.Sprintf("[%s] %s %s %s\n", tm.Timestamp(), msg.Topic(), tm.Event, tm.Desc)

		if _, err = f.WriteString(lineToWrite); err != nil {
			// TODO replace with logfile
			panic(err)
		}
	}
}

func (l *Listener) Start() {
	options := mqtt.NewClientOptions().AddBroker(l.Url)
	options.SetClientID("eventr")
	options.SetDefaultPublishHandler(l.MessageHandler())

	client := mqtt.NewClient(options)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	defer client.Disconnect(250)

	if token := client.Subscribe(l.TopicName, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	//Publish 5 messages to the topic at qos 1 and wait for the receipt
	//from the server after sending each message
	for i := 0; i < 5; i++ {
		fmt.Println("Publish sample message")
		text := sampleMessage()
		token := client.Publish("owntracks/phil/iPhone/event", 0, false, text)
		token.Wait()
		time.Sleep(500 * time.Millisecond)
	}

	//unsubscribe from topic
	if token := client.Unsubscribe(l.TopicName); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func sampleMessage() string {
	obj := &TransitionMessage{
		Tst:   time.Now().Unix(),
		Event: "enter",
		Desc:  "Test Desc"}
	msg, _ := json.Marshal(obj)
	return string(msg)
}
