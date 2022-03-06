package champ

import (
	math "github.com/chewxy/math32"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cbase"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"
)

type TrajectoryPlanner struct {
	leg                  *cbase.QuadLeg
	previousFootPosition cstructs.Transformation

	totalControlPoints int
	factorial          [13]float32
	refControlPointsX  [12]float32
	refControlPointsY  [12]float32
	controlPointsX     [12]float32
	controlPointsY     [12]float32

	heightRatio float32
	lengthRatio float32

	runOnce bool
}

func (trajectoryPlanner *TrajectoryPlanner) Init(leg *cbase.QuadLeg) *TrajectoryPlanner {

	*trajectoryPlanner = TrajectoryPlanner{
		leg:                  leg,
		previousFootPosition: cstructs.Transformation{},
		totalControlPoints:   12,
		factorial: [13]float32{
			1.0, 1.0, 2.0, 6.0, 24.0,
			120.0, 720.0, 5040.0, 40320.0,
			362880.0, 3628800.0,
			39916800.0, 479001600.0},
		refControlPointsX: [12]float32{-0.15, -0.2805, -0.3, -0.3, -0.3, 0.0, 0.0, 0.0, 0.3032, 0.3032, 0.2826, 0.15},
		refControlPointsY: [12]float32{-0.5, -0.5, -0.3611, -0.3611, -0.3611, -0.3611, -0.3611, -0.3214, -0.3214,
			-0.3214, -0.5, -0.5},
		controlPointsX: [12]float32{},
		controlPointsY: [12]float32{},
		heightRatio:    0,
		lengthRatio:    0,
		runOnce:        false,
	}

	return trajectoryPlanner
}

func (trajectoryPlanner *TrajectoryPlanner) UpdateControlPointsHeight(swingHeight float32) {
	newHeightRatio := swingHeight / 0.15

	if trajectoryPlanner.heightRatio != newHeightRatio {
		trajectoryPlanner.heightRatio = newHeightRatio

		for i := 0; i < 12; i++ {
			trajectoryPlanner.controlPointsY[i] = -((trajectoryPlanner.refControlPointsY[i] * trajectoryPlanner.heightRatio) + (0.5 * trajectoryPlanner.heightRatio))
		}
	}
}

func (trajectoryPlanner *TrajectoryPlanner) UpdateControlPointsLength(stepLength float32) {
	newLengthRatio := stepLength / 0.4

	if trajectoryPlanner.lengthRatio != newLengthRatio {
		trajectoryPlanner.lengthRatio = newLengthRatio

		for i := 0; i < 12; i++ {
			if i == 0 {
				trajectoryPlanner.controlPointsX[i] = -stepLength / 2.0
			} else if i == 11 {
				trajectoryPlanner.controlPointsX[i] = stepLength / 2.0
			} else {
				trajectoryPlanner.controlPointsX[i] = trajectoryPlanner.refControlPointsX[i] * trajectoryPlanner.lengthRatio
			}
		}
	}
}

func (trajectoryPlanner *TrajectoryPlanner) Generate(footPosition cstructs.Transformation, stepLength float32,
	rotation float32, swingPhaseSignal float32, stancePhaseSignal float32) cstructs.Transformation {

	trajectoryPlanner.UpdateControlPointsHeight(trajectoryPlanner.leg.GaitConfig().SwingHeight)

	if !trajectoryPlanner.runOnce {
		trajectoryPlanner.runOnce = true
		trajectoryPlanner.previousFootPosition = footPosition
	}

	if stepLength == 0 {
		trajectoryPlanner.previousFootPosition = footPosition
		trajectoryPlanner.leg.SetGaitPhase(true)
		return footPosition
	}

	trajectoryPlanner.UpdateControlPointsLength(stepLength)

	n := trajectoryPlanner.totalControlPoints - 1
	var x float32
	var y float32

	if stancePhaseSignal > swingPhaseSignal {
		trajectoryPlanner.leg.SetGaitPhase(true)

		x = (stepLength / 2) * (1 - (2 * stancePhaseSignal))
		y = -trajectoryPlanner.leg.GaitConfig().StanceDepth * math.Cos((math.Pi*x)/stepLength)
	} else if stancePhaseSignal < swingPhaseSignal {
		trajectoryPlanner.leg.SetGaitPhase(false)

		for i := 0; i < trajectoryPlanner.totalControlPoints; i++ {
			coeff := trajectoryPlanner.factorial[n] / (trajectoryPlanner.factorial[i] * trajectoryPlanner.factorial[n-i])

			x = x + coeff*math.Pow(swingPhaseSignal, float32(i))*math.Pow(1-swingPhaseSignal, float32(n-i))*trajectoryPlanner.controlPointsX[i]
			y = y - coeff*math.Pow(swingPhaseSignal, float32(i))*math.Pow(1-swingPhaseSignal, float32(n-i))*trajectoryPlanner.controlPointsY[i]
		}
	}

	footPosition.SetX(footPosition.X() + (x * math.Cos(rotation)))
	footPosition.SetY(footPosition.Y() + (x * math.Sin(rotation)))
	footPosition.SetZ(footPosition.Z() + y)

	if (swingPhaseSignal == 0.0 && stancePhaseSignal == 0.0) && stepLength > 0.0 {
		footPosition = trajectoryPlanner.previousFootPosition
	}

	trajectoryPlanner.previousFootPosition = footPosition

	return footPosition

}
