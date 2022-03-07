package champ

import (
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cbase"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"
)

type BodyController struct {
	quadBase *cbase.QuadBase
}

func (bodyController *BodyController) Init(quadBase *cbase.QuadBase) *BodyController {
	*bodyController = BodyController{
		quadBase: quadBase,
	}

	return bodyController
}

func (bodyController *BodyController) PoseCommand(footPositions [4]cstructs.Transformation, pose *cstructs.Pose) [4]cstructs.Transformation {

	for i := 0; i < 4; i++ {
		footPositions[i] = bodyController.poseCommandF(footPositions[i], bodyController.quadBase.Legs[i], pose)
	}

	return footPositions
}

func (bodyController *BodyController) poseCommandF(footPosition cstructs.Transformation, quadLeg *cbase.QuadLeg,
	pose *cstructs.Pose) cstructs.Transformation {

	reqTranslationX := -pose.Position.X()
	reqTranslationY := -pose.Position.Y()
	reqTranslationZ := -(quadLeg.ZeroStance().Z() + pose.Position.Z())
	maxTranslationZ := -quadLeg.ZeroStance().Z() * 0.65

	if reqTranslationZ < 0.0 {
		reqTranslationZ = 0.0
	} else if reqTranslationZ > maxTranslationZ {
		reqTranslationZ = maxTranslationZ
	}

	footPosition = quadLeg.ZeroStance()
	footPosition = footPosition.Translate(reqTranslationX, reqTranslationY, reqTranslationZ)

	footPosition = footPosition.RotateZ(-pose.Orientation.Z())
	footPosition = footPosition.RotateY(-pose.Orientation.Y())
	footPosition = footPosition.RotateX(-pose.Orientation.X())

	footPosition = KinematicsTransformToHip(footPosition, quadLeg)

	return footPosition
}
