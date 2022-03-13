package managers

import (
	"github.com/r4stl1n/micro-hal/code/internal/controller-node/handlers"
	"github.com/r4stl1n/micro-hal/code/pkg/consts"
	"github.com/r4stl1n/micro-hal/code/pkg/messages"
	"github.com/r4stl1n/micro-hal/code/pkg/mq"
	"github.com/r4stl1n/micro-hal/code/pkg/structs"
	"github.com/sirupsen/logrus"
)

type NodeManager struct {
	nats *mq.Nats

	poseHandler *handlers.PoseHandler
}

func (nodeManager *NodeManager) Init() *NodeManager {
	*nodeManager = NodeManager{
		nats: new(mq.Nats).Init(*new(structs.NatsConfig).Defaults()),
	}

	nodeManager.poseHandler = new(handlers.PoseHandler).Init(nodeManager.nats)

	return nodeManager
}

func (nodeManager *NodeManager) connectToNats() error {
	return nodeManager.nats.Connect()
}

func (nodeManager *NodeManager) Process() error {

	connectToNatsError := nodeManager.connectToNats()
	if connectToNatsError != nil {
		return connectToNatsError
	}

	receiveChannel := make(chan *[]byte, 100)
	_, bindError := nodeManager.nats.EncodedConn.BindRecvChan(consts.MQPoseSetChannel, receiveChannel)
	if bindError != nil {
		return bindError
	}

	logrus.Info("service started waiting for messages")

	for {
		select {
		case receiveData := <-receiveChannel:
			requestMessage := new(messages.Message)
			requestError := requestMessage.Unpack(*receiveData)
			if requestError != nil {
				logrus.Error(requestError)
				continue
			}

			logrus.Debugf("Received Message: %+v", requestMessage)

			switch requestMessage.Type {

			case messages.PoseMessage:
				sMessage := new(messages.Pose)

				unpackError := sMessage.Unpack(requestMessage.Data)
				if unpackError != nil {
					logrus.Error(unpackError)
					continue
				}

				nodeManager.poseHandler.Handle(sMessage)

			default:
				logrus.Errorf("unknown message received %v+", requestMessage)
			}

		}
	}
}
