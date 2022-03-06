package hmath

import (
	"fmt"

	"github.com/barnex/fmath"
)

type Mat3 [9]float32

func (mat3 *Mat3) Pointer() *[9]float32 { return (*[9]float32)(mat3) }
func (mat3 *Mat3) Slice() []float32     { return mat3[:] }

func (mat3 *Mat3) String() string {
	return fmt.Sprintf("[%f,%f,%f,\n %f,%f,%f,\n %f,%f,%f,]", mat3[0], mat3[3], mat3[6], mat3[1], mat3[4], mat3[7], mat3[2], mat3[5], mat3[8])
}

//mat3[0] mat3[1] mat3[2]
//mat3[3] mat3[4] mat3[5]
//mat3[6] mat3[7] mat3[8]

func (mat3 *Mat3) SetAt(row int, column int, value float32) {
	selection := 0

	switch row {
	case 0:
		selection = 0
	case 1:
		selection = 3
	case 2:
		selection = 6
	}

	selection = selection + column

	mat3[selection] = value
}

func (mat3 *Mat3) GetAt(row int, column int) float32 {
	selection := 0

	switch row {
	case 0:
		selection = 0
	case 1:
		selection = 3
	case 2:
		selection = 6
	}

	selection = selection + column

	return mat3[selection]
}

func Mat3Identity() Mat3 {
	return Mat3{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1}
}

func Mat3Translate(v Vec2) Mat3 {
	return Mat3{
		1, 0, 0,
		0, 1, 0,
		v[0], v[1], 1}
}

func Mat3Scale(v Vec2) Mat3 {
	return Mat3{
		v[0], 0, 0,
		0, v[1], 0,
		0, 0, 1}
}

func Mat3Rotate(radians float32) Mat3 {
	s, c := fmath.Sincos(radians)
	return Mat3{
		c, s, 0,
		-s, c, 0,
		0, 0, 1}
}

func (mat3 Mat3) Mul(m2 Mat3) Mat3 {
	return Mat3{
		mat3[0]*m2[0] + mat3[1]*m2[3] + mat3[2]*m2[6],
		mat3[0]*m2[1] + mat3[1]*m2[4] + mat3[2]*m2[7],
		mat3[0]*m2[2] + mat3[1]*m2[5] + mat3[2]*m2[8],
		mat3[3]*m2[0] + mat3[4]*m2[3] + mat3[5]*m2[6],
		mat3[3]*m2[1] + mat3[4]*m2[4] + mat3[5]*m2[7],
		mat3[3]*m2[2] + mat3[4]*m2[5] + mat3[5]*m2[8],
		mat3[6]*m2[0] + mat3[7]*m2[3] + mat3[8]*m2[6],
		mat3[6]*m2[1] + mat3[7]*m2[4] + mat3[8]*m2[7],
		mat3[6]*m2[2] + mat3[7]*m2[5] + mat3[8]*m2[8]}
}

func (mat3 Mat3) Invert() Mat3 {
	identity := 1.0 / (mat3[0]*mat3[4]*mat3[8] + mat3[3]*mat3[7]*mat3[2] + mat3[6]*mat3[1]*mat3[5] - mat3[6]*mat3[4]*mat3[2] - mat3[3]*mat3[1]*mat3[8] - mat3[0]*mat3[7]*mat3[5])

	return Mat3{
		(mat3[4]*mat3[8] - mat3[5]*mat3[7]) * identity,
		(mat3[2]*mat3[7] - mat3[1]*mat3[8]) * identity,
		(mat3[1]*mat3[5] - mat3[2]*mat3[4]) * identity,
		(mat3[5]*mat3[6] - mat3[3]*mat3[8]) * identity,
		(mat3[0]*mat3[8] - mat3[2]*mat3[6]) * identity,
		(mat3[2]*mat3[3] - mat3[0]*mat3[5]) * identity,
		(mat3[3]*mat3[7] - mat3[4]*mat3[6]) * identity,
		(mat3[1]*mat3[6] - mat3[0]*mat3[7]) * identity,
		(mat3[0]*mat3[4] - mat3[1]*mat3[3]) * identity}
}
