package messages

import "github.com/vmihailenco/msgpack/v5"

type Pose struct {
	X     float32
	Y     float32
	Z     float32
	Roll  float32
	Pitch float32
	Yaw   float32
}

func (pose *Pose) Init() *Pose {
	*pose = Pose{}
	return pose
}

func (pose *Pose) Pack() []byte {
	bytes, _ := msgpack.Marshal(&pose)
	return bytes
}

func (pose *Pose) Unpack(data []byte) error {
	return msgpack.Unmarshal(data, &pose)
}
