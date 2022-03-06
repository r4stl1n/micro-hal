package hmath

import (
	"fmt"
	math "github.com/chewxy/math32"

	"github.com/barnex/fmath"
)

type Vec3 [3]float32

func (vec3 *Vec3) Pointer() *[3]float32            { return (*[3]float32)(vec3) }
func (vec3 *Vec3) Slice() []float32                { return vec3[:] }
func (vec3 Vec3) X() float32                       { return vec3[0] }
func (vec3 Vec3) Y() float32                       { return vec3[1] }
func (vec3 Vec3) Z() float32                       { return vec3[2] }
func (vec3 Vec3) XY() (float32, float32)           { return vec3[0], vec3[1] }
func (vec3 Vec3) XZ() (float32, float32)           { return vec3[0], vec3[2] }
func (vec3 Vec3) YX() (float32, float32)           { return vec3[1], vec3[0] }
func (vec3 Vec3) YZ() (float32, float32)           { return vec3[1], vec3[2] }
func (vec3 Vec3) ZX() (float32, float32)           { return vec3[2], vec3[0] }
func (vec3 Vec3) ZY() (float32, float32)           { return vec3[2], vec3[1] }
func (vec3 Vec3) XYZ() (float32, float32, float32) { return vec3[0], vec3[1], vec3[2] }
func (vec3 Vec3) XYVec() Vec2                      { return Vec2{vec3[0], vec3[1]} }
func (vec3 Vec3) XZVec() Vec2                      { return Vec2{vec3[0], vec3[2]} }
func (vec3 Vec3) YXVec() Vec2                      { return Vec2{vec3[1], vec3[0]} }
func (vec3 Vec3) YZVec() Vec2                      { return Vec2{vec3[1], vec3[2]} }
func (vec3 Vec3) ZXVec() Vec2                      { return Vec2{vec3[2], vec3[0]} }
func (vec3 Vec3) ZYVec() Vec2                      { return Vec2{vec3[2], vec3[1]} }

func (vec3 *Vec3) SetX(x float32) { vec3[0] = x }
func (vec3 *Vec3) SetY(y float32) { vec3[1] = y }
func (vec3 *Vec3) SetZ(z float32) { vec3[2] = z }
func (vec3 *Vec3) SetXYZ(x, y, z float32) {
	vec3[0] = x
	vec3[1] = y
	vec3[2] = z
	return
}

func (vec3 Vec3) String() string {
	return fmt.Sprintf("[%f,%f,%f]", vec3[0], vec3[1], vec3[2])
}

func (vec3 Vec3) Add(v2 Vec3) Vec3 {
	return Vec3{vec3[0] + v2[0], vec3[1] + v2[1], vec3[2] + v2[2]}
}

func (vec3 Vec3) Sub(v2 Vec3) Vec3 {
	return Vec3{vec3[0] - v2[0], vec3[1] - v2[1], vec3[2] - v2[2]}
}

func (vec3 Vec3) Mul(v2 Vec3) Vec3 {
	return Vec3{vec3[0] * v2[0], vec3[1] * v2[1], vec3[2] * v2[2]}
}

func (vec3 Vec3) MulF(f float32) Vec3 {
	return Vec3{vec3[0] * f, vec3[1] * f, vec3[2] * f}
}

func (vec3 Vec3) MulMat4(m Mat4) Vec3 {
	vec3, w := vec3.MulMat4W(m)
	return vec3.DivF(w)
}

func (vec3 Vec3) MulMat4W(m Mat4) (Vec3, float32) {
	return Vec3{
			m[0]*vec3[0] + m[4]*vec3[1] + m[8]*vec3[2] + m[12],
			m[1]*vec3[0] + m[5]*vec3[1] + m[9]*vec3[2] + m[13],
			m[2]*vec3[0] + m[6]*vec3[1] + m[10]*vec3[2] + m[14]},
		m[3]*vec3[0] + m[7]*vec3[1] + m[11]*vec3[2] + m[15]
}

func (vec3 Vec3) MulQuaternion(q Quaternion) Vec3 {
	return q.MulQuaternion(Quaternion{vec3[0], vec3[1], vec3[2], 0}).MulQuaternion(q.Conjugate()).XYZVec()
}

func (vec3 Vec3) Div(v2 Vec3) Vec3 {
	return Vec3{vec3[0] / v2[0], vec3[1] / v2[1], vec3[2] / v2[2]}
}

func (vec3 Vec3) DivF(f float32) Vec3 {
	return Vec3{vec3[0] / f, vec3[1] / f, vec3[2] / f}
}

func (vec3 Vec3) Dot(v2 Vec3) float32 {
	return vec3[0]*v2[0] + vec3[1]*v2[1] + vec3[2]*v2[2]
}

func (vec3 Vec3) Cross(v2 Vec3) Vec3 {
	return Vec3{
		vec3[1]*v2[2] - vec3[2]*v2[1],
		vec3[2]*v2[0] - vec3[0]*v2[2],
		vec3[0]*v2[1] - vec3[1]*v2[0]}
}

func (vec3 Vec3) Magnitude() float32 {
	var total float32

	for i := 0; i < 3; i++ {
		total = total + math.Pow(vec3[i], 2)
	}

	return math.Sqrt(total)
}

func (vec3 Vec3) Len() float32 {
	return fmath.Sqrt(vec3[0]*vec3[0] + vec3[1]*vec3[1] + vec3[2]*vec3[2])
}

func (vec3 Vec3) LenSqr() float32 {
	return vec3[0]*vec3[0] + vec3[1]*vec3[1] + vec3[2]*vec3[2]
}

func (vec3 Vec3) Norm() Vec3 {
	return vec3.MulF(1.0 / vec3.Len())
}

func (vec3 Vec3) AngleTo(v2 Vec3) float32 {
	return fmath.Acos(vec3.Norm().Dot(v2.Norm()))
}
