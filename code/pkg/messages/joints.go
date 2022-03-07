package messages

import (
	"github.com/r4stl1n/micro-hal/code/pkg/hmath"
	"github.com/vmihailenco/msgpack/v5"
)

type Joints struct {
	LeftFront  hmath.Vec3
	RightFront hmath.Vec3
	LeftBack   hmath.Vec3
	RightBack  hmath.Vec3
}

func (joints *Joints) Init() *Joints {
	*joints = Joints{}
	return joints
}

func (joints *Joints) Pack() []byte {
	bytes, _ := msgpack.Marshal(&joints)
	return bytes
}

func (joints *Joints) Unpack(data []byte) error {
	return msgpack.Unmarshal(data, &joints)
}
