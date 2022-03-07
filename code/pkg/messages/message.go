package messages

import (
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
)

type MessageType int

const (
	ExampleRequestMessage  MessageType = 1
	ExampleResponseMessage MessageType = 2

	JointsMessage MessageType = 3
	PoseMessage   MessageType = 4
)

type Message struct {
	Type     MessageType
	Data     []byte
	RespChan string
	Multi    bool
}

func (message *Message) Request(respChan string) *Message {
	*message = Message{
		RespChan: respChan,
		Multi:    false,
	}

	return message
}

func (message *Message) Response() *Message {
	*message = Message{}

	return message
}

func (message *Message) ResponseMulti(multi bool) *Message {
	*message = Message{
		Multi: multi,
	}

	return message
}

func (message *Message) Pack() []byte {
	bytes, _ := msgpack.Marshal(&message)
	return bytes
}

func (message *Message) Unpack(data []byte) error {
	return msgpack.Unmarshal(data, &message)
}

func (message *Message) packMessage(response interface{}) {
	switch response.(type) {

	case *ExampleRequest:
		message.Type = ExampleRequestMessage
		message.Data = response.(*ExampleRequest).Pack()
	case *ExampleResponse:
		message.Type = ExampleResponseMessage
		message.Data = response.(*ExampleResponse).Pack()
	case *Joints:
		message.Type = JointsMessage
		message.Data = response.(*Joints).Pack()
	case *Pose:
		message.Type = PoseMessage
		message.Data = response.(*Pose).Pack()

	default:
		logrus.Errorf("Unknown message type %v+", response)
	}
}

func (message *Message) Build(response interface{}) []byte {

	message.packMessage(response)

	return message.Pack()
}

func (message *Message) BuildM(response interface{}) *Message {

	message.packMessage(response)

	return message
}
