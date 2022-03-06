package hmath

import (
	"fmt"

	"github.com/barnex/fmath"
)

type Mat4 [16]float32

func (mat4 *Mat4) Pointer() *[16]float32 { return (*[16]float32)(mat4) }
func (mat4 *Mat4) Slice() []float32      { return mat4[:] }

func (mat4 *Mat4) String() string {
	return fmt.Sprintf("[%f,%f,%f,%f,\n %f,%f,%f,%f,\n %f,%f,%f,%f,\n %f,%f,%f,%f]",
		mat4[0], mat4[4], mat4[8], mat4[12], mat4[1], mat4[5], mat4[9], mat4[13], mat4[2], mat4[6], mat4[10], mat4[14], mat4[3], mat4[7], mat4[11], mat4[15])
}

func (mat4 *Mat4) SetAt(row int, column int, value float32) {
	selection := 0

	switch row {
	case 0:
		selection = 0
	case 1:
		selection = 4
	case 2:
		selection = 8
	case 3:
		selection = 12
	}

	selection = selection + column

	mat4[selection] = value
}

func (mat4 *Mat4) GetAt(row int, column int) float32 {
	selection := 0

	switch row {
	case 0:
		selection = 0
	case 1:
		selection = 4
	case 2:
		selection = 8
	case 3:
		selection = 12
	}

	selection = selection + column

	return mat4[selection]
}

func Mat4Identity() Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1}
}

func Mat4Translate(v Vec3) Mat4 {
	return Mat4{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		v[0], v[1], v[2], 1}
}

func Mat4Scale(v Vec3) Mat4 {
	return Mat4{
		v[0], 0, 0, 0,
		0, v[1], 0, 0,
		0, 0, v[2], 0,
		0, 0, 0, 1}
}

func Mat4Rotate(axis Vec3, radians float32) Mat4 {
	x, y, z := axis.Norm().XYZ()
	c := fmath.Cos(radians)
	s := fmath.Sin(radians)

	return Mat4{
		x*x*(1-c) + c, y*x*(1-c) + z*s, x*z*(1-c) - y*s, 0,
		x*y*(1-c) - z*s, y*y*(1-c) + c, y*z*(1-c) + x*s, 0,
		x*z*(1-c) + y*s, y*z*(1-c) - x*s, z*z*(1-c) + c, 0,
		0, 0, 0, 1}
}

func Mat4LookAt(eye, center, up Vec3) Mat4 {
	f := center.Sub(eye).Norm()
	s := f.Cross(up.Norm())
	u := s.Cross(f)
	return Mat4{
		s[0], u[0], -f[0], 0,
		s[1], u[1], -f[1], 0,
		s[2], u[2], -f[2], 0,
		-s.Dot(eye), -u.Dot(eye), f.Dot(eye), 1}
}

func Mat4Frustum(left, right, bottom, top, zNear, zFar float32) Mat4 {
	width := right - left
	height := top - bottom
	depth := zFar - zNear
	return Mat4{
		(zNear * 2.0) / width, 0, 0, 0,
		0, (zNear * 2.0) / height, 0, 0,
		(left + right) / width, (bottom + top) / height, -(zNear + zFar) / depth, -1,
		0, 0, -(zNear * zFar * 2.0) / depth, 0}
}

func Mat4Perspective(fovY, aspect, zNear, zFar float32) Mat4 {
	f := 1.0 / fmath.Tan(fovY/2.0)
	d := zNear - zFar
	return Mat4{
		f / aspect, 0, 0, 0,
		0, f, 0, 0,
		0, 0, (zFar + zNear) / d, -1,
		0, 0, (2 * zFar * zNear) / d, 0}
}

func (mat4 Mat4) Mul(m2 Mat4) Mat4 {
	return Mat4{
		mat4[0]*m2[0] + mat4[1]*m2[4] + mat4[2]*m2[8] + mat4[3]*m2[12],
		mat4[0]*m2[1] + mat4[1]*m2[5] + mat4[2]*m2[9] + mat4[3]*m2[13],
		mat4[0]*m2[2] + mat4[1]*m2[6] + mat4[2]*m2[10] + mat4[3]*m2[14],
		mat4[0]*m2[3] + mat4[1]*m2[7] + mat4[2]*m2[11] + mat4[3]*m2[15],
		mat4[4]*m2[0] + mat4[5]*m2[4] + mat4[6]*m2[8] + mat4[7]*m2[12],
		mat4[4]*m2[1] + mat4[5]*m2[5] + mat4[6]*m2[9] + mat4[7]*m2[13],
		mat4[4]*m2[2] + mat4[5]*m2[6] + mat4[6]*m2[10] + mat4[7]*m2[14],
		mat4[4]*m2[3] + mat4[5]*m2[7] + mat4[6]*m2[11] + mat4[7]*m2[15],
		mat4[8]*m2[0] + mat4[9]*m2[4] + mat4[10]*m2[8] + mat4[11]*m2[12],
		mat4[8]*m2[1] + mat4[9]*m2[5] + mat4[10]*m2[9] + mat4[11]*m2[13],
		mat4[8]*m2[2] + mat4[9]*m2[6] + mat4[10]*m2[10] + mat4[11]*m2[14],
		mat4[8]*m2[3] + mat4[9]*m2[7] + mat4[10]*m2[11] + mat4[11]*m2[15],
		mat4[12]*m2[0] + mat4[13]*m2[4] + mat4[14]*m2[8] + mat4[15]*m2[12],
		mat4[12]*m2[1] + mat4[13]*m2[5] + mat4[14]*m2[9] + mat4[15]*m2[13],
		mat4[12]*m2[2] + mat4[13]*m2[6] + mat4[14]*m2[10] + mat4[15]*m2[14],
		mat4[12]*m2[3] + mat4[13]*m2[7] + mat4[14]*m2[11] + mat4[15]*m2[15]}
}

func (mat4 Mat4) Invert() Mat4 {
	var s, c [6]float32
	s[0] = mat4[0]*mat4[5] - mat4[4]*mat4[1]
	s[1] = mat4[0]*mat4[6] - mat4[4]*mat4[2]
	s[2] = mat4[0]*mat4[7] - mat4[4]*mat4[3]
	s[3] = mat4[1]*mat4[6] - mat4[5]*mat4[2]
	s[4] = mat4[1]*mat4[7] - mat4[5]*mat4[3]
	s[5] = mat4[2]*mat4[7] - mat4[6]*mat4[3]

	c[0] = mat4[8]*mat4[13] - mat4[12]*mat4[9]
	c[1] = mat4[8]*mat4[14] - mat4[12]*mat4[10]
	c[2] = mat4[8]*mat4[15] - mat4[12]*mat4[11]
	c[3] = mat4[9]*mat4[14] - mat4[13]*mat4[10]
	c[4] = mat4[9]*mat4[15] - mat4[13]*mat4[11]
	c[5] = mat4[10]*mat4[15] - mat4[14]*mat4[11]

	// assumes it is invertible
	identity := 1.0 / (s[0]*c[5] - s[1]*c[4] + s[2]*c[3] + s[3]*c[2] - s[4]*c[1] + s[5]*c[0])

	return Mat4{
		(mat4[5]*c[5] - mat4[6]*c[4] + mat4[7]*c[3]) * identity,
		(-mat4[1]*c[5] + mat4[2]*c[4] - mat4[3]*c[3]) * identity,
		(mat4[13]*s[5] - mat4[14]*s[4] + mat4[15]*s[3]) * identity,
		(-mat4[9]*s[5] + mat4[10]*s[4] - mat4[11]*s[3]) * identity,
		(-mat4[4]*c[5] + mat4[6]*c[2] - mat4[7]*c[1]) * identity,
		(mat4[0]*c[5] - mat4[2]*c[2] + mat4[3]*c[1]) * identity,
		(-mat4[12]*s[5] + mat4[14]*s[2] - mat4[15]*s[1]) * identity,
		(mat4[8]*s[5] - mat4[10]*s[2] + mat4[11]*s[1]) * identity,
		(mat4[4]*c[4] - mat4[5]*c[2] + mat4[7]*c[0]) * identity,
		(-mat4[0]*c[4] + mat4[1]*c[2] - mat4[3]*c[0]) * identity,
		(mat4[12]*s[4] - mat4[13]*s[2] + mat4[15]*s[0]) * identity,
		(-mat4[8]*s[4] + mat4[9]*s[2] - mat4[11]*s[0]) * identity,
		(-mat4[4]*c[3] + mat4[5]*c[1] - mat4[6]*c[0]) * identity,
		(mat4[0]*c[3] - mat4[1]*c[1] + mat4[2]*c[0]) * identity,
		(-mat4[12]*s[3] + mat4[13]*s[1] - mat4[14]*s[0]) * identity,
		(mat4[8]*s[3] - mat4[9]*s[1] + mat4[10]*s[0]) * identity}
}
