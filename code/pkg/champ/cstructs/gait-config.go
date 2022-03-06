package cstructs

import "github.com/r4stl1n/micro-hal/code/pkg/hmath"

type GaitConfig struct {
	KneeOrientation   string
	PantographLeg     bool
	OdomScalar        float32
	MaxLinearVelocity hmath.Vec3
	ComXTranslation   float32
	SwingHeight       float32
	StanceDepth       float32
	StanceDuration    float32
	NominalHeight     float32
}

func (gaitConfig *GaitConfig) Init(kneeOrientation string, pantoLeg bool, odomScalar float32,
	maxLinearVelocity hmath.Vec3, comXTrans float32, swingHeight float32, stanceDepth float32,
	stanceDuration float32, nominalHeight float32) *GaitConfig {

	*gaitConfig = GaitConfig{
		KneeOrientation:   kneeOrientation,
		PantographLeg:     pantoLeg,
		OdomScalar:        odomScalar,
		MaxLinearVelocity: maxLinearVelocity,
		ComXTranslation:   comXTrans,
		SwingHeight:       swingHeight,
		StanceDepth:       stanceDepth,
		StanceDuration:    stanceDuration,
		NominalHeight:     nominalHeight,
	}

	return gaitConfig
}
