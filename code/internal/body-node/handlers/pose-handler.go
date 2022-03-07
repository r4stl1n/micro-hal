package handlers

import (
	"github.com/r4stl1n/micro-hal/code/pkg/messages"
	"github.com/r4stl1n/micro-hal/code/pkg/mq"
)

type PoseHandler struct {
	nats *mq.Nats
}

func (poseHandler *PoseHandler) Init(nats *mq.Nats) *PoseHandler {
	*poseHandler = PoseHandler{
		nats: nats,
	}

	return poseHandler
}

func (poseHandler *PoseHandler) Handle(message *messages.Pose) {

}
