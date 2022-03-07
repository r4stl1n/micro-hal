package mq

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/r4stl1n/micro-hal/code/pkg/consts"
	"github.com/r4stl1n/micro-hal/code/pkg/messages"
	"github.com/r4stl1n/micro-hal/code/pkg/structs"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

type Nats struct {
	Config      structs.NatsConfig
	Conn        *nats.Conn
	EncodedConn *nats.EncodedConn
	Error       error
}

func (n *Nats) Init(cfg structs.NatsConfig) *Nats {
	*n = Nats{
		Config: cfg,
	}

	return n
}

func (n *Nats) Connect() error {
	if n.Error != nil {
		return n.Error
	}

	descriptor := fmt.Sprintf("mq://%s:%s@%s", n.Config.User, n.Config.Pass, n.Config.Host)

	n.Conn, n.Error = nats.Connect(descriptor, nats.ReconnectWait(time.Second))
	if n.Error != nil {
		return n.Error
	}

	n.EncodedConn, n.Error = nats.NewEncodedConn(n.Conn, nats.DEFAULT_ENCODER)

	return n.Error
}

func (n *Nats) SendAwaitResponse(channel string, message *messages.Message) (*messages.Message, error) {
	newResponseChannel := consts.MQNodePrefix + uuid.NewV4().String()

	receiveChannel := make(chan *[]byte, 100)
	subscription, bindError := n.EncodedConn.BindRecvChan(newResponseChannel, receiveChannel)
	if bindError != nil {
		return nil, bindError
	}

	defer func() {
		unsubscribeError := subscription.Unsubscribe()

		if unsubscribeError != nil {
			logrus.Error(unsubscribeError)
		}
	}()

	message.RespChan = newResponseChannel

	publishError := n.EncodedConn.Publish(channel, message.Pack())
	if publishError != nil {
		return nil, publishError
	}

	select {
	case receiveData := <-receiveChannel:
		unpackMessage := new(messages.Message)
		unpackError := unpackMessage.Unpack(*receiveData)

		if unpackError != nil {
			return nil, unpackError
		}

		return unpackMessage, nil

	case <-time.After(5 * time.Second):
		return nil, fmt.Errorf("did not receive response for requested message")
	}
}

func (n *Nats) SendAwaitResponseMultiSingle(channel string, message *messages.Message) ([]*messages.Message, error) {
	newResponseChannel := consts.MQNodePrefix + uuid.NewV4().String()

	receiveChannel := make(chan *[]byte, 100)
	subscription, bindError := n.EncodedConn.BindRecvChan(newResponseChannel, receiveChannel)
	if bindError != nil {
		return nil, bindError
	}

	defer func() {
		unsubscribeError := subscription.Unsubscribe()
		if unsubscribeError != nil {
			logrus.Error(unsubscribeError)
		}
	}()

	message.RespChan = newResponseChannel

	publishError := n.EncodedConn.Publish(channel, message.Pack())
	if publishError != nil {
		return nil, publishError
	}

	t := time.NewTimer(5 * time.Second)
	results := make([]*messages.Message, 0)

	for {
		select {
		case receiveData := <-receiveChannel:
			unpackMessage := new(messages.Message)
			unpackError := unpackMessage.Unpack(*receiveData)
			if unpackError != nil {
				return nil, unpackError
			}

			results = append(results, unpackMessage)
			if !unpackMessage.Multi {
				return results, nil
			}

			t.Reset(5 * time.Second)

		case <-t.C:
			return nil, fmt.Errorf("did not receive response for requested message")
		}
	}
}

func (n *Nats) SendAwaitResponseMulti(channel string, message *messages.Message) ([]*messages.Message, error) {
	newResponseChannel := consts.MQNodePrefix + uuid.NewV4().String()

	receiveChannel := make(chan *[]byte, 100)
	subscription, bindError := n.EncodedConn.BindRecvChan(newResponseChannel, receiveChannel)
	if bindError != nil {
		return nil, bindError
	}

	defer func() {
		unsubscribeError := subscription.Unsubscribe()
		if unsubscribeError != nil {
			logrus.Error(unsubscribeError)
		}
	}()

	message.RespChan = newResponseChannel

	publishError := n.EncodedConn.Publish(channel, message.Pack())
	if publishError != nil {
		return nil, publishError
	}

	results := make([]*messages.Message, 0)

	for {
		select {
		case receiveData := <-receiveChannel:
			unpackMessage := new(messages.Message)
			unpackError := unpackMessage.Unpack(*receiveData)
			if unpackError != nil {
				return nil, unpackError
			}

			results = append(results, unpackMessage)

		case <-time.After(5 * time.Second):
			return results, nil
		}
	}
}
