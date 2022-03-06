package hmath

import (
	"fmt"

	"github.com/barnex/fmath"
)

type Quaternion [4]float32

func (quaternion Quaternion) String() string {
	return fmt.Sprintf("[%f,%f,%f,%f]", quaternion[0], quaternion[1], quaternion[2], quaternion[3])
}

func (quaternion Quaternion) XYZVec() Vec3 {
	return Vec3{quaternion[0], quaternion[1], quaternion[2]}
}

func QuaternionAxisRotation(axis Vec3, angle float32) Quaternion {
	var d = 1 / axis.Len()
	halfAngle := angle / 2.0
	s := fmath.Sin(halfAngle)
	c := fmath.Cos(halfAngle)
	return Quaternion{s * axis[0] * d, s * axis[1] * d, s * axis[2] * d, c}
}

func QuaternionIdentity() Quaternion {
	return Quaternion{0, 0, 0, 1}
}

func QuaternionPitchYawRoll(pitch, yaw, roll float32) Quaternion {
	yawS, yawC := fmath.Sincos(yaw * 0.5)
	pitchS, pitchC := fmath.Sincos(pitch * 0.5)
	rollS, rollC := fmath.Sincos(roll * 0.5)
	var q Quaternion
	q[0] = pitchS*yawC*rollC - pitchC*yawS*rollS
	q[1] = pitchC*yawS*rollC + pitchS*yawC*rollS
	q[2] = pitchC*yawC*rollS - pitchS*yawS*rollC
	q[3] = pitchC*yawC*rollC + pitchS*yawS*rollS
	return q
}

func (quaternion Quaternion) MulQuaternion(q2 Quaternion) Quaternion {
	return Quaternion{
		quaternion[3]*q2[0] + quaternion[0]*q2[3] + quaternion[1]*q2[2] - quaternion[2]*q2[1],
		quaternion[3]*q2[1] + quaternion[1]*q2[3] + quaternion[2]*q2[0] - quaternion[0]*q2[2],
		quaternion[3]*q2[2] + quaternion[2]*q2[3] + quaternion[0]*q2[1] - quaternion[1]*q2[0],
		quaternion[3]*q2[3] - quaternion[0]*q2[0] - quaternion[1]*q2[1] - quaternion[2]*q2[2]}
}

func (quaternion Quaternion) Conjugate() Quaternion {
	return Quaternion{-quaternion[0], -quaternion[1], -quaternion[2], quaternion[3]}
}

func (quaternion Quaternion) AngleAround(axis Vec3) Quaternion {
	return Quaternion{-quaternion[0], -quaternion[1], -quaternion[2], quaternion[3]}
}

func (quaternion Quaternion) Mat4() Mat4 {
	x, y, z, w := quaternion[0], quaternion[1], quaternion[2], quaternion[3]
	return Mat4{
		1 - 2*y*y - 2*z*z, 2*x*y + 2*w*z, 2*x*z - 2*w*y, 0,
		2*x*y - 2*w*z, 1 - 2*x*x - 2*z*z, 2*y*z + 2*w*x, 0,
		2*x*z + 2*w*y, 2*y*z - 2*w*x, 1 - 2*x*x - 2*y*y, 0,
		0, 0, 0, 1,
	}
}
