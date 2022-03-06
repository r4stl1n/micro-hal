package hmath

import (
	"fmt"

	"github.com/barnex/fmath"
)

type Vec2 [2]float32

func (vec2 *Vec2) Pointer() *[2]float32 { return (*[2]float32)(vec2) }
func (vec2 *Vec2) Slice() []float32     { return vec2[:] }
func (vec2 Vec2) X() float32            { return vec2[0] }
func (vec2 Vec2) Y() float32            { return vec2[1] }
func (vec2 Vec2) XY() (x, y float32) {
	x = vec2[0]
	y = vec2[1]
	return
}

func (vec2 *Vec2) SetX(x float32) { vec2[0] = x }
func (vec2 *Vec2) SetY(y float32) { vec2[1] = y }
func (vec2 *Vec2) SetXY(x, y float32) {
	vec2[0] = x
	vec2[1] = y
	return
}

func (vec2 Vec2) String() string {
	return fmt.Sprintf("[%f,%f]", vec2[0], vec2[1])
}

func (vec2 Vec2) Add(v2 Vec2) Vec2 {
	return Vec2{vec2[0] + v2[0], vec2[1] + v2[1]}
}

func (vec2 Vec2) Sub(v2 Vec2) Vec2 {
	return Vec2{vec2[0] - v2[0], vec2[1] - v2[1]}
}

func (vec2 Vec2) Mul(v2 Vec2) Vec2 {
	return Vec2{vec2[0] * v2[0], vec2[1] * v2[1]}
}

func (vec2 Vec2) MulF(f float32) Vec2 {
	return Vec2{vec2[0] * f, vec2[1] * f}
}

func (vec2 Vec2) MulMat3(m Mat3) (Vec2, float32) {
	return Vec2{
			m[0]*vec2[0] + m[3]*vec2[1] + m[6],
			m[1]*vec2[0] + m[4]*vec2[1] + m[7]},
		m[2]*vec2[0] + m[5]*vec2[1] + m[8]
}

func (vec2 Vec2) Div(v2 Vec2) Vec2 {
	return Vec2{vec2[0] / v2[0], vec2[1] / v2[1]}
}

func (vec2 Vec2) DivF(f float32) Vec2 {
	return Vec2{vec2[0] / f, vec2[1] / f}
}

func (vec2 Vec2) Dot(v2 Vec2) float32 {
	return vec2[0]*v2[0] + vec2[1]*v2[1]
}

func (vec2 Vec2) Len() float32 {
	return fmath.Sqrt(vec2[0]*vec2[0] + vec2[1]*vec2[1])
}

func (vec2 Vec2) LenSqr() float32 {
	return vec2[0]*vec2[0] + vec2[1]*vec2[1]
}

func (vec2 Vec2) Norm() Vec2 {
	return vec2.MulF(1.0 / vec2.Len())
}

func (vec2 Vec2) Atan2() float32 {
	return fmath.Atan2(vec2[1], vec2[0])
}

func (vec2 Vec2) AngleTo(v2 Vec2) float32 {
	return fmath.Acos(vec2.Norm().Dot(v2.Norm()))
}
