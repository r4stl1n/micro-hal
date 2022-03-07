package champ

import (
	math "github.com/chewxy/math32"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cbase"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"
)

type Kinematics struct {
	quadBase *cbase.QuadBase
}

func (kinematics *Kinematics) Init(quadBase *cbase.QuadBase) *Kinematics {
	*kinematics = Kinematics{
		quadBase: quadBase,
	}

	return kinematics
}

func (kinematics *Kinematics) Inverse(jointPositions [12]float32, footPositions [4]cstructs.Transformation) [12]float32 {

	calculatedJoints := [12]float32{}

	for i := 0; i < 4; i++ {

		calculatedJoints[(i * 3)], calculatedJoints[(i*3)+1], calculatedJoints[(i*3)+2] = kinematics.inverseF(kinematics.quadBase.Legs[i], footPositions[i])

		if math.IsNaN(calculatedJoints[(i*3)]) || math.IsNaN(calculatedJoints[(i*3)+1]) || math.IsNaN(calculatedJoints[(i*3)+2]) {
			return jointPositions
		}
	}

	for i := 0; i < 12; i++ {
		jointPositions[i] = calculatedJoints[i]
	}

	return jointPositions
}

func (kinematics *Kinematics) inverseF(quadLeg *cbase.QuadLeg, footPosition cstructs.Transformation) (float32, float32, float32) {

	hipJoint := float32(0.0)
	lowerLegJoint := float32(0.0)
	upperLegJoint := float32(0.0)

	tempFootPosition := footPosition

	l0 := float32(0.0)

	for i := 1; i < 4; i++ {
		l0 = l0 + quadLeg.JointChain[i].Y()
	}

	l1 := -math.Sqrt(math.Pow(quadLeg.LowerLegJoint.X(), 2) + math.Pow(quadLeg.LowerLegJoint.Z(), 2))
	ikAlpha := math.Acos(quadLeg.LowerLegJoint.X()/l1) - (math.Pi / 2)

	l2 := -math.Sqrt(math.Pow(quadLeg.FootJoint.X(), 2) + math.Pow(quadLeg.FootJoint.Z(), 2))
	ikBeta := math.Acos(quadLeg.FootJoint.X()/l2) - (math.Pi / 2)

	x := tempFootPosition.X()
	y := tempFootPosition.Y()
	z := tempFootPosition.Z()

	hipJoint = -(math.Atan(y/z) - ((math.Pi / 2) - math.Acos(-l0/math.Sqrt(math.Pow(y, 2)+math.Pow(z, 2)))))
	tempFootPosition = tempFootPosition.RotateX(-hipJoint)
	tempFootPosition = tempFootPosition.Translate(-quadLeg.UpperLegJoint.X(), 0.0, -quadLeg.UpperLegJoint.Z())

	x = tempFootPosition.X()
	y = tempFootPosition.Y()
	z = tempFootPosition.Z()

	targetToFoot := math.Sqrt(math.Pow(x, 2) + math.Pow(z, 2))

	if targetToFoot >= (math.Abs(l1) + math.Abs(l2)) {
		return hipJoint, lowerLegJoint, upperLegJoint
	}

	lowerLegJoint = float32(quadLeg.KneeDirection()) * math.Acos((math.Pow(z, 2)+math.Pow(x, 2)-math.Pow(l1, 2)-math.Pow(l2, 2))/(2*l1*l2))
	upperLegJoint = math.Atan(x/z) - math.Atan((l2*math.Sin(lowerLegJoint))/(l1+(l2*math.Cos(lowerLegJoint))))
	lowerLegJoint = lowerLegJoint + (ikBeta - ikAlpha)
	upperLegJoint = upperLegJoint + ikAlpha

	if quadLeg.KneeDirection() < 0 {
		if upperLegJoint < 0 {
			upperLegJoint = upperLegJoint + math.Pi
		}
	} else {
		if upperLegJoint > 0 {
			upperLegJoint = upperLegJoint + math.Pi
		}
	}

	return hipJoint, lowerLegJoint, upperLegJoint
}

func KinematicsTransformToHip(footPosition cstructs.Transformation, quadLeg *cbase.QuadLeg) cstructs.Transformation {
	return footPosition.Translate(-quadLeg.HipJoint.X(), -quadLeg.HipJoint.Y(), -quadLeg.HipJoint.Z())
}

func KinematicsTransformToBase(footPosition cstructs.Transformation, quadLeg *cbase.QuadLeg) cstructs.Transformation {
	return footPosition.Translate(quadLeg.HipJoint.X(), quadLeg.HipJoint.Y(), quadLeg.HipJoint.Z())
}
