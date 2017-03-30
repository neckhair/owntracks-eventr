package listener

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/neckhair/owntracks-eventr/utils"
)

type Listener struct {
	client    mqtt.Client
	Config    *Configuration
	TopicName string
	TLSConfig *tls.Config
}

func NewListener(config *Configuration) *Listener {
	return &Listener{
		TopicName: "owntracks/+/+/event",
		Config:    config,
		TLSConfig: &tls.Config{},
	}
}

func (l *Listener) MessageHandler() mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		tm, err := NewTransitionMessage(msg.Payload())

		if err != nil {
			utils.Error("Failure while parsing payload.")
			utils.Error(string(msg.Payload()))
			return
		}

		f, err := os.OpenFile(l.Config.Filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			utils.Error("Cannot open output file to write event data.")
			return
		}
		defer f.Close()

		lineToWrite := fmt.Sprintf("[%s] %s %s %s\n", tm.Timestamp(), msg.Topic(), tm.Event, tm.Desc)

		if _, err = f.WriteString(lineToWrite); err != nil {
			utils.Error("Error writing event data to file.")
			return
		}
	}
}

func (l *Listener) Start() error {
	l.client = mqtt.NewClient(l.ClientOptions())
	if token := l.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	if token := l.client.Subscribe(l.TopicName, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// Unsubscribe from topic
func (l *Listener) Stop() {
	log.Println("Disconnecting from broker.")
	if token := l.client.Unsubscribe(l.TopicName); token.Wait() && token.Error() != nil {
		utils.Error("Could not unsubscribe from MQTT topic.")
		log.Println(token.Error())
	}
	l.client.Disconnect(250)
}

func (l *Listener) ClientOptions() *mqtt.ClientOptions {
	options := mqtt.NewClientOptions().AddBroker(l.Config.Url)
	options.SetClientID(l.clientID())
	options.AutoReconnect = true
	options.SetDefaultPublishHandler(l.MessageHandler())
	options.SetTLSConfig(l.TLSConfig)
	options.Username = l.Config.Username
	options.Password = l.Config.Password

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

func (l *Listener) clientID() string {
	pid := os.Getpid()
	return "eventr-" + strconv.Itoa(pid)
}

func sampleMessage() string {
	obj := &TransitionMessage{
		Tst:   time.Now().Unix(),
		Event: "enter",
		Desc:  "Test Desc"}
	msg, _ := json.Marshal(obj)
	return string(msg)
}
