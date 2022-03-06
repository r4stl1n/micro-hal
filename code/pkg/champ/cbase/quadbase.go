package cbase

import "github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"

type QuadBase struct {
	speed cstructs.Velocities

	legs []*QuadLeg

	leftFront  *QuadLeg
	rightFront *QuadLeg
	leftBack   *QuadLeg
	rightBack  *QuadLeg

	gaitConfig cstructs.GaitConfig
}

func (quadBase *QuadBase) Init(gaitConfig cstructs.GaitConfig) *QuadBase {

	*quadBase = QuadBase{
		gaitConfig: gaitConfig,
		leftFront:  new(QuadLeg).Init(),
		rightFront: new(QuadLeg).Init(),
		leftBack:   new(QuadLeg).Init(),
		rightBack:  new(QuadLeg).Init(),
	}

	quadBase.legs = append(quadBase.legs, quadBase.leftFront)
	quadBase.legs = append(quadBase.legs, quadBase.rightFront)
	quadBase.legs = append(quadBase.legs, quadBase.leftBack)
	quadBase.legs = append(quadBase.legs, quadBase.rightBack)

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
		retVal = append(retVal, quadBase.legs[i].hipJoint.Theta())
		retVal = append(retVal, quadBase.legs[i].upperLegJoint.Theta())
		retVal = append(retVal, quadBase.legs[i].lowerLegJoint.Theta())
	}

	return retVal
}

func (quadBase *QuadBase) GetFootPositions() []float32 {
	var retVal []float32

	for i := 0; i < 4; i++ {
		retVal = append(retVal, quadBase.legs[i].footJoint.Theta())
	}

	return retVal
}

func (quadBase *QuadBase) UpdateJointPositions(positions []float32) {
	for i := 0; i < 4; i++ {
		index := i * 3

		quadBase.legs[i].hipJoint.SetTheta(positions[index])
		quadBase.legs[i].upperLegJoint.SetTheta(positions[index+1])
		quadBase.legs[i].lowerLegJoint.SetTheta(positions[index+2])
	}
}

func (quadBase *QuadBase) SetGaitConfig(gaitConfig cstructs.GaitConfig) {
	quadBase.gaitConfig = gaitConfig

	for i := 0; i < 4; i++ {
		dir := 0

		quadBase.legs[i].SetId(i)

		if i < 2 {
			dir = quadBase.GetKneeDirection(gaitConfig.KneeOrientation[0:1])
		} else {
			dir = quadBase.GetKneeDirection(gaitConfig.KneeOrientation[1:])
		}

		quadBase.legs[i].SetPantograph(gaitConfig.PantographLeg)
		quadBase.legs[i].SetKneeDirection(dir)
		quadBase.legs[i].SetGaitConfig(gaitConfig)
	}
}

func (quadBase *QuadBase) GaitConfig() cstructs.GaitConfig {
	return quadBase.gaitConfig
}
