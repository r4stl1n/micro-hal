package champ

import (
	math "github.com/chewxy/math32"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cbase"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"
	"time"
)

type LegController struct {
	quadBase *cbase.QuadBase

	trajectoryPlanners [4]*TrajectoryPlanner

	phaseGenerator *PhaseGenerator

	leftFrontTrajectoryPlanner  *TrajectoryPlanner
	rightFrontTrajectoryPlanner *TrajectoryPlanner
	leftBackTrajectoryPlanner   *TrajectoryPlanner
	rightBackTrajectoryPlanner  *TrajectoryPlanner
}

func (legController *LegController) Init(quadBase *cbase.QuadBase, currentTime time.Time) *LegController {

	*legController = LegController{
		quadBase:                    quadBase,
		trajectoryPlanners:          [4]*TrajectoryPlanner{},
		phaseGenerator:              new(PhaseGenerator).Init(quadBase, currentTime),
		leftFrontTrajectoryPlanner:  new(TrajectoryPlanner).Init(quadBase.LeftFront),
		rightFrontTrajectoryPlanner: new(TrajectoryPlanner).Init(quadBase.RightFront),
		leftBackTrajectoryPlanner:   new(TrajectoryPlanner).Init(quadBase.LeftBack),
		rightBackTrajectoryPlanner:  new(TrajectoryPlanner).Init(quadBase.RightBack),
	}

	legController.trajectoryPlanners[0] = legController.leftFrontTrajectoryPlanner
	legController.trajectoryPlanners[1] = legController.rightFrontTrajectoryPlanner
	legController.trajectoryPlanners[2] = legController.leftBackTrajectoryPlanner
	legController.trajectoryPlanners[3] = legController.rightBackTrajectoryPlanner

	return legController
}

func (legController *LegController) capVelocities(velocity float32, minVelocity float32, maxVelocity float32) float32 {

	if velocity < minVelocity {
		return minVelocity
	}

	if velocity > maxVelocity {
		return maxVelocity
	}

	return velocity
}

func (legController *LegController) TransformLeg(quadLeg *cbase.QuadLeg, stepLength float32, rotation float32,
	stepX float32, stepY float32, theta float32) (float32, float32) {

	transformedStance := quadLeg.ZeroStance()
	transformedStance.Translate(stepX, stepY, 0.0)
	transformedStance.RotateZ(theta)

	zeroStance := quadLeg.ZeroStance()

	deltaX := transformedStance.X() - zeroStance.X()
	deltaY := transformedStance.Y() - zeroStance.Y()

	stepLength = math.Sqrt(math.Pow(deltaX, 2)+math.Pow(deltaY, 2)) * 2.0

	rotation = math.Atan2(deltaY, deltaX)

	return stepLength, rotation

}

func (legController *LegController) RaibertHeuristic(stanceDuration float32, targetVelocity float32) float32 {
	return (stanceDuration / 2.0) * targetVelocity
}

func (legController *LegController) VelocityCommand(footPositions [4]cstructs.Transformation, velocities cstructs.Velocities,
	currentTime time.Time) ([4]cstructs.Transformation, cstructs.Velocities) {

	velocities.Linear.SetX(
		legController.capVelocities(
			velocities.Linear.X(),
			-legController.quadBase.GaitConfig().MaxLinearVelocity.X(),
			legController.quadBase.GaitConfig().MaxLinearVelocity.X()))

	velocities.Linear.SetY(
		legController.capVelocities(
			velocities.Linear.Y(),
			-legController.quadBase.GaitConfig().MaxLinearVelocity.Y(),
			legController.quadBase.GaitConfig().MaxLinearVelocity.Y()))

	velocities.Angular.SetZ(
		legController.capVelocities(
			velocities.Angular.Z(),
			-legController.quadBase.GaitConfig().MaxAngularVelocity,
			legController.quadBase.GaitConfig().MaxAngularVelocity))

	tangentialVelocity := velocities.Angular.Z() * legController.quadBase.LeftFront.CenterToNominal()
	velocity := math.Sqrt(math.Pow(velocities.Linear.X(), 2) + math.Pow(velocities.Linear.Y()+tangentialVelocity, 2))

	stepX := legController.RaibertHeuristic(legController.quadBase.GaitConfig().StanceDuration, velocities.Linear.X())
	stepY := legController.RaibertHeuristic(legController.quadBase.GaitConfig().StanceDuration, velocities.Linear.Y())
	stepTheta := legController.RaibertHeuristic(legController.quadBase.GaitConfig().StanceDuration, tangentialVelocity)

	theta := math.Sin((stepTheta/2)/legController.quadBase.LeftFront.CenterToNominal()) * 2

	stepLengths := [4]float32{0.0, 0.0, 0.0, 0.0}
	trajectoryRotations := [4]float32{0.0, 0.0, 0.0, 0.0}
	sumOfSteps := float32(0.0)

	for i := 0; i < 4; i++ {
		legController.TransformLeg(legController.quadBase.Legs[i], stepLengths[i], trajectoryRotations[i], stepX, stepY, theta)
		sumOfSteps = sumOfSteps + stepLengths[i]
	}

	legController.phaseGenerator.Run(velocity, sumOfSteps/4.0, currentTime)

	for i := 0; i < 4; i++ {
		legController.trajectoryPlanners[i].Generate(footPositions[i], stepLengths[i], trajectoryRotations[i],
			legController.phaseGenerator.swingPhaseSignal[i], legController.phaseGenerator.stancePhaseSignal[i])
	}

	return footPositions, velocities
}
