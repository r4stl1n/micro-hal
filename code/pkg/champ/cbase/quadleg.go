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

	HipJoint      *QuadJoint
	UpperLegJoint *QuadJoint
	LowerLegJoint *QuadJoint
	FootJoint     *QuadJoint

	gaitConfig cstructs.GaitConfig

	JointChain []*QuadJoint
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
		HipJoint:        new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),
		UpperLegJoint:   new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),
		LowerLegJoint:   new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),
		FootJoint:       new(QuadJoint).Init(hmath.Vec3{}, hmath.Vec3{}, 0.0),

		gaitConfig: cstructs.GaitConfig{},
	}

	quadLeg.JointChain = append(quadLeg.JointChain, quadLeg.HipJoint)
	quadLeg.JointChain = append(quadLeg.JointChain, quadLeg.UpperLegJoint)
	quadLeg.JointChain = append(quadLeg.JointChain, quadLeg.LowerLegJoint)
	quadLeg.JointChain = append(quadLeg.JointChain, quadLeg.FootJoint)

	return quadLeg
}

func (quadLeg *QuadLeg) FootFromHip() cstructs.Transformation {

	var footPosition cstructs.Transformation

	for i := 3; i > 0; i-- {
		footPosition = footPosition.Translate(quadLeg.JointChain[i].X(), quadLeg.JointChain[i].Y(), quadLeg.JointChain[i].Z())

		if i > 1 {
			footPosition = footPosition.RotateY(quadLeg.JointChain[i-1].Theta())
		}

	}

	return footPosition
}

func (quadLeg *QuadLeg) FootFromBase() cstructs.Transformation {
	var footPosition cstructs.Transformation

	footPosition.Point = quadLeg.FootFromHip().Point
	footPosition = footPosition.RotateX(quadLeg.HipJoint.Theta())
	footPosition = footPosition.Translate(quadLeg.HipJoint.X(), quadLeg.HipJoint.Y(), quadLeg.HipJoint.Z())

	return footPosition
}

func (quadLeg *QuadLeg) SetJoints(hipJoint float32, upperLegJoint float32, lowerLegJoint float32) {
	quadLeg.HipJoint.SetTheta(hipJoint)
	quadLeg.UpperLegJoint.SetTheta(upperLegJoint)
	quadLeg.LowerLegJoint.SetTheta(lowerLegJoint)
}

func (quadLeg QuadLeg) ZeroStance() cstructs.Transformation {
	quadLeg.zeroStance.SetX(quadLeg.HipJoint.X() + quadLeg.UpperLegJoint.X() + quadLeg.gaitConfig.ComXTranslation)
	quadLeg.zeroStance.SetY(quadLeg.HipJoint.Y() + quadLeg.UpperLegJoint.Y())
	quadLeg.zeroStance.SetZ(quadLeg.HipJoint.Z() + quadLeg.UpperLegJoint.Z() + quadLeg.LowerLegJoint.Z() + quadLeg.FootJoint.Z())

	return quadLeg.zeroStance
}

func (quadLeg *QuadLeg) CenterToNominal() float32 {
	x := quadLeg.HipJoint.X() + quadLeg.UpperLegJoint.X()
	y := quadLeg.HipJoint.Y() + quadLeg.UpperLegJoint.Y()

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
