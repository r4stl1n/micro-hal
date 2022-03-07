package cbase

import "github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"

type QuadBase struct {
	speed cstructs.Velocities

	Legs []*QuadLeg

	LeftFront  *QuadLeg
	RightFront *QuadLeg
	LeftBack   *QuadLeg
	RightBack  *QuadLeg

	gaitConfig cstructs.GaitConfig
}

func (quadBase *QuadBase) Init(gaitConfig cstructs.GaitConfig) *QuadBase {

	*quadBase = QuadBase{
		gaitConfig: gaitConfig,
		LeftFront:  new(QuadLeg).Init(),
		RightFront: new(QuadLeg).Init(),
		LeftBack:   new(QuadLeg).Init(),
		RightBack:  new(QuadLeg).Init(),
	}

	quadBase.Legs = append(quadBase.Legs, quadBase.LeftFront)
	quadBase.Legs = append(quadBase.Legs, quadBase.RightFront)
	quadBase.Legs = append(quadBase.Legs, quadBase.LeftBack)
	quadBase.Legs = append(quadBase.Legs, quadBase.RightBack)

	return quadBase
}

func (quadBase *QuadBase) GetKneeDirection(character string) int {
	switch character {
	case ">":
		return -1
	case "<":
		return 1
	default:
		return -1
	}
}

func (quadBase *QuadBase) GetJointPositions() []float32 {
	var retVal []float32

	for i := 0; i < 4; i++ {
		retVal = append(retVal, quadBase.Legs[i].HipJoint.Theta())
		retVal = append(retVal, quadBase.Legs[i].UpperLegJoint.Theta())
		retVal = append(retVal, quadBase.Legs[i].LowerLegJoint.Theta())
	}

	return retVal
}

func (quadBase *QuadBase) GetFootPositions() []float32 {
	var retVal []float32

	for i := 0; i < 4; i++ {
		retVal = append(retVal, quadBase.Legs[i].FootJoint.Theta())
	}

	return retVal
}

func (quadBase *QuadBase) UpdateJointPositions(positions []float32) {
	for i := 0; i < 4; i++ {
		index := i * 3

		quadBase.Legs[i].HipJoint.SetTheta(positions[index])
		quadBase.Legs[i].UpperLegJoint.SetTheta(positions[index+1])
		quadBase.Legs[i].LowerLegJoint.SetTheta(positions[index+2])
	}
}

func (quadBase *QuadBase) SetGaitConfig(gaitConfig cstructs.GaitConfig) {
	quadBase.gaitConfig = gaitConfig

	for i := 0; i < 4; i++ {
		dir := 0

		quadBase.Legs[i].SetId(i)

		if i < 2 {
			dir = quadBase.GetKneeDirection(gaitConfig.KneeOrientation[0:1])
		} else {
			dir = quadBase.GetKneeDirection(gaitConfig.KneeOrientation[1:])
		}

		quadBase.Legs[i].SetPantograph(gaitConfig.PantographLeg)
		quadBase.Legs[i].SetKneeDirection(dir)
		quadBase.Legs[i].SetGaitConfig(gaitConfig)
	}
}

func (quadBase *QuadBase) GaitConfig() cstructs.GaitConfig {
	return quadBase.gaitConfig
}
