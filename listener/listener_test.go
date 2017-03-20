package listener_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	l "github.com/neckhair/owntracks-eventr/listener"
)

func TestClientOptions(t *testing.T) {
	config := &l.Configuration{
		Url:      "tcp://localhost:8883",
		Username: "stinky",
		Password: "supersecret",
	}
	listener := l.NewListener(config)

	clientOptions := listener.ClientOptions()
	assert.Equal(t, "stinky", clientOptions.Username)
	assert.Equal(t, "supersecret", clientOptions.Password)
	assert.Equal(t, "eventr", clientOptions.ClientID)
	assert.Equal(t, true, clientOptions.AutoReconnect)
	assert.Equal(t, "localhost:8883", clientOptions.Servers[0].Host)
}
