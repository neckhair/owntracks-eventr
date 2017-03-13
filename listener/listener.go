package listener

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Listener struct {
	client         mqtt.Client
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
	l.client = mqtt.NewClient(l.ClientOptions())
	if token := l.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := l.client.Subscribe(l.TopicName, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

// Unsubscribe from topic
func (l *Listener) Stop() {
	if token := l.client.Unsubscribe(l.TopicName); token.Wait() && token.Error() != nil {
		// TODO replace with log
		fmt.Println(token.Error())
	}
	l.client.Disconnect(250)
}

func (l *Listener) ClientOptions() *mqtt.ClientOptions {
	options := mqtt.NewClientOptions().AddBroker(l.Url)
	options.SetClientID("eventr")
	options.SetDefaultPublishHandler(l.MessageHandler())
	return options
}

//Publish sample messages to the topic at qos 1 and wait for the receipt
//from the server after sending each message
func (l *Listener) PublishExampleMessages(number int) {
	for i := 0; i < number; i++ {
		text := sampleMessage()
		token := l.client.Publish("owntracks/phil/iPhone/event", 0, false, text)
		token.Wait()
		time.Sleep(500 * time.Millisecond)
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
