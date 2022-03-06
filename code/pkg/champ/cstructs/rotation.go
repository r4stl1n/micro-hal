package cstructs

import (
	math "github.com/chewxy/math32"
	"github.com/r4stl1n/micro-hal/code/pkg/hmath"
)

type Rotation struct {
	Mat3 hmath.Mat3
}

func (rotation *Rotation) FromEulerAngles(psi float32, theta float32, phi float32) *Rotation {

	*rotation = Rotation{}

	rotation.Mat3.SetAt(0, 0, math.Cos(psi)*math.Cos(theta))
	rotation.Mat3.SetAt(1, 0, math.Cos(theta)*math.Sin(phi))
	rotation.Mat3.SetAt(2, 0, -math.Sin(theta))

	rotation.Mat3.SetAt(0, 1, math.Cos(phi)*math.Sin(psi)*math.Sin(theta)-math.Cos(psi)*math.Sin(phi))
	rotation.Mat3.SetAt(1, 1, math.Cos(psi)*math.Cos(phi)+math.Sin(psi)*math.Sin(phi)*math.Sin(theta))
	rotation.Mat3.SetAt(2, 1, math.Cos(psi)*math.Cos(theta))

	rotation.Mat3.SetAt(0, 2, math.Sin(psi)*math.Sin(phi)+math.Cos(psi)*math.Cos(phi)*math.Sin(theta))
	rotation.Mat3.SetAt(1, 2, math.Cos(psi)*math.Sin(phi)*math.Sin(theta)-math.Cos(phi)*math.Sin(psi))
	rotation.Mat3.SetAt(2, 2, math.Cos(psi)*math.Cos(theta))

	return rotation
}

func (rotation *Rotation) RotateX(phi float32) *Rotation {

	var temp1 float32
	var temp2 float32

	temp1 = rotation.Mat3.GetAt(1, 0)*math.Cos(phi) - rotation.Mat3.GetAt(2, 0)*math.Sin(phi)
	temp2 = rotation.Mat3.GetAt(2, 0)*math.Cos(phi) + rotation.Mat3.GetAt(1, 0)*math.Sin(phi)
	rotation.Mat3.SetAt(1, 0, temp1)
	rotation.Mat3.SetAt(2, 0, temp2)

	temp1 = rotation.Mat3.GetAt(1, 1)*math.Cos(phi) - rotation.Mat3.GetAt(2, 1)*math.Sin(phi)
	temp2 = rotation.Mat3.GetAt(2, 1)*math.Cos(phi) + rotation.Mat3.GetAt(1, 1)*math.Sin(phi)
	rotation.Mat3.SetAt(1, 1, temp1)
	rotation.Mat3.SetAt(2, 1, temp2)

	temp1 = rotation.Mat3.GetAt(1, 2)*math.Cos(phi) - rotation.Mat3.GetAt(2, 2)*math.Sin(phi)
	temp2 = rotation.Mat3.GetAt(2, 2)*math.Cos(phi) + rotation.Mat3.GetAt(1, 2)*math.Sin(phi)
	rotation.Mat3.SetAt(1, 2, temp1)
	rotation.Mat3.SetAt(2, 2, temp2)

	return rotation
}

func (rotation *Rotation) RotateY(theta float32) *Rotation {

	var temp1 float32
	var temp2 float32

	temp1 = rotation.Mat3.GetAt(0, 0)*math.Cos(theta) + rotation.Mat3.GetAt(2, 0)*math.Sin(theta)
	temp2 = rotation.Mat3.GetAt(2, 0)*math.Cos(theta) - rotation.Mat3.GetAt(0, 0)*math.Sin(theta)
	rotation.Mat3.SetAt(0, 0, temp1)
	rotation.Mat3.SetAt(2, 0, temp2)

	temp1 = rotation.Mat3.GetAt(0, 1)*math.Cos(theta) + rotation.Mat3.GetAt(2, 1)*math.Sin(theta)
	temp2 = rotation.Mat3.GetAt(2, 1)*math.Cos(theta) - rotation.Mat3.GetAt(0, 1)*math.Sin(theta)
	rotation.Mat3.SetAt(0, 1, temp1)
	rotation.Mat3.SetAt(2, 1, temp2)

	temp1 = rotation.Mat3.GetAt(0, 2)*math.Cos(theta) + rotation.Mat3.GetAt(2, 2)*math.Sin(theta)
	temp2 = rotation.Mat3.GetAt(2, 2)*math.Cos(theta) - rotation.Mat3.GetAt(0, 2)*math.Sin(theta)
	rotation.Mat3.SetAt(1, 2, temp1)
	rotation.Mat3.SetAt(1, 2, temp2)

	return rotation
}

func (rotation *Rotation) RotateZ(psi float32) *Rotation {

	var temp1 float32
	var temp2 float32

	temp1 = rotation.Mat3.GetAt(0, 0)*math.Cos(psi) - rotation.Mat3.GetAt(1, 0)*math.Sin(psi)
	temp2 = rotation.Mat3.GetAt(1, 0)*math.Cos(psi) + rotation.Mat3.GetAt(0, 0)*math.Sin(psi)
	rotation.Mat3.SetAt(0, 0, temp1)
	rotation.Mat3.SetAt(1, 0, temp2)

	temp1 = rotation.Mat3.GetAt(0, 1)*math.Cos(psi) - rotation.Mat3.GetAt(1, 1)*math.Sin(psi)
	temp2 = rotation.Mat3.GetAt(1, 1)*math.Cos(psi) + rotation.Mat3.GetAt(0, 1)*math.Sin(psi)
	rotation.Mat3.SetAt(0, 1, temp1)
	rotation.Mat3.SetAt(1, 1, temp2)

	temp1 = rotation.Mat3.GetAt(0, 2)*math.Cos(psi) - rotation.Mat3.GetAt(1, 2)*math.Sin(psi)
	temp2 = rotation.Mat3.GetAt(1, 2)*math.Cos(psi) + rotation.Mat3.GetAt(0, 2)*math.Sin(psi)
	rotation.Mat3.SetAt(0, 2, temp1)
	rotation.Mat3.SetAt(1, 2, temp2)

	return rotation
}
