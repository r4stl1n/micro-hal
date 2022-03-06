package cbase

import (
	math "github.com/chewxy/math32"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"
	"github.com/r4stl1n/micro-hal/code/pkg/hmath"
)

type QuadLeg struct {
	numberOfLinks   int
	zeroStance      cstructs.Transformation
	centerToNominal float32
	id              int

	lastTouchDown float32

	inContact     bool
	kneeDirection int
	isPantograph  bool
	gaitPhase     bool

	hipJoint      *QuadJoint
	upperLegJoint *QuadJoint
	lowerLegJoint *QuadJoint
	footJoint     *QuadJoint

	gaitConfig cstructs.GaitConfig

	jointChain []*QuadJoint
}

func (quadLeg *QuadLeg) Init() *QuadLeg {
	*quadLeg = QuadLeg{
		numberOfLinks:   4,
		zeroStance:      cstructs.Transformation{},
		centerToNominal: 0,
		id:              0,
		lastTouchDown:   0,
		inContact:       true,
		kneeDirection:   0,
		isPantograph:    false,
		gaitPhase:       true,
		hipJoint:        new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),
		upperLegJoint:   new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),
		lowerLegJoint:   new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),
		footJoint:       new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),

		gaitConfig: cstructs.GaitConfig{},
	}

	quadLeg.jointChain = append(quadLeg.jointChain, quadLeg.hipJoint)
	quadLeg.jointChain = append(quadLeg.jointChain, quadLeg.upperLegJoint)
	quadLeg.jointChain = append(quadLeg.jointChain, quadLeg.lowerLegJoint)
	quadLeg.jointChain = append(quadLeg.jointChain, quadLeg.footJoint)

	return quadLeg
}

func (quadLeg *QuadLeg) FootFromHip() cstructs.Transformation {

	var footPosition cstructs.Transformation

	for i := 3; i > 0; i-- {
		footPosition.Translate(quadLeg.jointChain[i].X(), quadLeg.jointChain[i].Y(), quadLeg.jointChain[i].Z())

		if i > 1 {
			footPosition.RotateY(quadLeg.jointChain[i-1].Theta())
		}

	}

	return footPosition
}

func (quadLeg *QuadLeg) FootFromBase() cstructs.Transformation {
	var footPosition cstructs.Transformation

	footPosition.Point = quadLeg.FootFromHip().Point
	footPosition.RotateX(quadLeg.hipJoint.Theta())
	footPosition.Translate(quadLeg.hipJoint.X(), quadLeg.hipJoint.Y(), quadLeg.hipJoint.Z())

	return footPosition
}

func (quadLeg *QuadLeg) SetJoints(hipJoint float32, upperLegJoint float32, lowerLegJoint float32) {
	quadLeg.hipJoint.SetTheta(hipJoint)
	quadLeg.upperLegJoint.SetTheta(upperLegJoint)
	quadLeg.lowerLegJoint.SetTheta(lowerLegJoint)
}

func (quadLeg *QuadLeg) ZeroStance() cstructs.Transformation {
	quadLeg.zeroStance.SetX(quadLeg.hipJoint.X() + quadLeg.upperLegJoint.X() + quadLeg.gaitConfig.ComXTranslation)
	quadLeg.zeroStance.SetY(quadLeg.hipJoint.Y() + quadLeg.upperLegJoint.Y())
	quadLeg.zeroStance.SetZ(quadLeg.hipJoint.Z() + quadLeg.upperLegJoint.Z() + quadLeg.lowerLegJoint.Z() + quadLeg.footJoint.Z())

	return quadLeg.zeroStance
}

func (quadLeg *QuadLeg) CenterToNominal() float32 {
	x := quadLeg.hipJoint.X() + quadLeg.upperLegJoint.X()
	y := quadLeg.hipJoint.Y() + quadLeg.upperLegJoint.Y()

	return math.Sqrt(math.Pow(x, 2) + math.Pow(y, 2))
}

func (quadLeg *QuadLeg) Id() int {
	return quadLeg.id
}

func (quadLeg *QuadLeg) SetInContact(contact bool) {
	quadLeg.inContact = contact
}

func (quadLeg *QuadLeg) IsInContact() bool {
	return quadLeg.inContact
}

func (quadLeg *QuadLeg) SetGaitPhase(gaitPhase bool) {
	quadLeg.gaitPhase = gaitPhase
}

func (quadLeg *QuadLeg) GaitConfig() cstructs.GaitConfig {
	return quadLeg.gaitConfig
}

func (quadLeg *QuadLeg) IsInGaitPhase() bool {
	return quadLeg.gaitPhase
}

func (quadLeg *QuadLeg) SetKneeDirection(kneeDirection int) {
	quadLeg.kneeDirection = kneeDirection
}

func (quadLeg *QuadLeg) KneeDirection() int {
	return quadLeg.kneeDirection
}

func (quadLeg *QuadLeg) SetPantograph(pantograph bool) {
	quadLeg.isPantograph = pantograph
}

func (quadLeg *QuadLeg) Pantograph() bool {
	return quadLeg.isPantograph
}

func (quadLeg *QuadLeg) SetGaitConfig(gaitConfig cstructs.GaitConfig) {
	quadLeg.gaitConfig = gaitConfig
}

func (quadLeg *QuadLeg) SetId(id int) {
	quadLeg.id = id
}
