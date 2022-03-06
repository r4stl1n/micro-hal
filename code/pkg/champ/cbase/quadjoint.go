package cbase

import (
	"github.com/r4stl1n/micro-hal/code/pkg/hmath"
)

type QuadJoint struct {
	theta       float32
	translation hmath.Vec3
	rotation    hmath.Vec3
}

func (quadJoint *QuadJoint) Init(point hmath.Vec3, euler hmath.Vec3, theta float32) *QuadJoint {
	*quadJoint = QuadJoint{
		theta:       theta,
		translation: point,
		rotation:    euler,
	}

	return quadJoint
}

func (quadJoint *QuadJoint) SetTheta(theta float32) {
	quadJoint.theta = theta
}

func (quadJoint *QuadJoint) SetTranslation(point hmath.Vec3) {
	quadJoint.translation = point
}

func (quadJoint *QuadJoint) SetRotation(euler hmath.Vec3) {
	quadJoint.rotation = euler
}

func (quadJoint *QuadJoint) SetOrigin(point hmath.Vec3, euler hmath.Vec3) {
	quadJoint.translation = point
	quadJoint.rotation = euler
}

func (quadJoint *QuadJoint) Theta() float32 {
	return quadJoint.theta
}

func (quadJoint *QuadJoint) X() float32 {
	return quadJoint.translation.X()
}

func (quadJoint *QuadJoint) Y() float32 {
	return quadJoint.translation.Y()
}

func (quadJoint *QuadJoint) Z() float32 {
	return quadJoint.translation.Z()
}

func (quadJoint *QuadJoint) Roll() float32 {
	return quadJoint.rotation.X()
}

func (quadJoint *QuadJoint) Pitch() float32 {
	return quadJoint.rotation.Y()
}

func (quadJoint *QuadJoint) Yaw() float32 {
	return quadJoint.rotation.Z()
}
