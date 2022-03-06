package cstructs

import (
	math "github.com/chewxy/math32"
	"github.com/r4stl1n/micro-hal/code/pkg/hmath"
)

type Transformation struct {
	Rotation Rotation
	Point    hmath.Vec3
}

func (transformation Transformation) X() float32 {
	return transformation.Point.X()
}

func (transformation Transformation) Y() float32 {
	return transformation.Point.Y()
}

func (transformation Transformation) Z() float32 {
	return transformation.Point.Z()
}

func (transformation Transformation) SetX(x float32) {
	transformation.Point.SetX(x)
}

func (transformation Transformation) SetY(y float32) {
	transformation.Point.SetY(y)
}

func (transformation Transformation) SetZ(z float32) {
	transformation.Point.SetZ(z)
}

func (transformation Transformation) Identify() Transformation {
	transformation.Rotation.Mat3 = hmath.Mat3Identity()
	transformation.Point = hmath.Vec3{}

	return transformation
}

func (transformation Transformation) RotateX(phi float32) Transformation {

	var temp hmath.Vec3
	transformation.Rotation.RotateX(phi)

	temp.SetX(transformation.Point.X())
	temp.SetY(math.Cos(phi)*transformation.Point.Y() - math.Sin(phi)*transformation.Point.Z())
	temp.SetZ(math.Sin(phi)*transformation.Point.Y() + math.Cos(phi)*transformation.Point.Z())

	transformation.Point = temp

	return transformation
}

func (transformation Transformation) RotateY(theta float32) Transformation {

	var temp hmath.Vec3
	transformation.Rotation.RotateY(theta)

	temp.SetX(math.Cos(theta)*transformation.Point.X() + math.Sin(theta)*transformation.Point.Z())
	temp.SetY(transformation.Point.Y())
	temp.SetZ(-math.Sin(theta)*transformation.Point.X() + math.Cos(theta)*transformation.Point.Z())

	transformation.Point = temp

	return transformation
}

func (transformation Transformation) RotateZ(psi float32) Transformation {

	var temp hmath.Vec3
	transformation.Rotation.RotateZ(psi)

	temp.SetX(math.Cos(psi)*transformation.Point.X() - math.Sin(psi)*transformation.Point.Y())
	temp.SetY(math.Sin(psi)*transformation.Point.X() + math.Cos(psi)*transformation.Point.Y())
	temp.SetZ(transformation.Point.Z())

	transformation.Point = temp

	return transformation
}

func (transformation Transformation) Translate(x float32, y float32, z float32) Transformation {
	transformation.Point.SetX(transformation.Point.X() + x)
	transformation.Point.SetY(transformation.Point.Y() + y)
	transformation.Point.SetZ(transformation.Point.Z() + z)

	return transformation
}
