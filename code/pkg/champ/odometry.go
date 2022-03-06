package champ

import (
	"github.com/chewxy/math32"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cbase"
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cstructs"
	"time"
)

type Odometry struct {
	quadBase              *cbase.QuadBase
	previousFootPositions [4]cstructs.Transformation

	previousFootContacts [4]bool
	previousTheta        [4]float32
	previousTime         time.Time
	previousVelocity     cstructs.Velocities
	beta                 float32
}

func (odometry *Odometry) Init(quadBase *cbase.QuadBase, currentTime time.Time) *Odometry {
	*odometry = Odometry{
		quadBase:             quadBase,
		previousFootContacts: [4]bool{true, true, true, true},
		previousTheta:        [4]float32{0.0, 0.0, 0.0, 0.0},
		previousTime:         currentTime,
		beta:                 0.1,
	}

	for i := 0; i < 4; i++ {
		odometry.previousFootPositions[i] = quadBase.Legs[i].FootFromBase()
	}

	return odometry
}

func (odometry *Odometry) allFeetInContact() bool {
	if odometry.quadBase.Legs[0].IsInContact() &&
		odometry.quadBase.Legs[1].IsInContact() &&
		odometry.quadBase.Legs[2].IsInContact() &&
		odometry.quadBase.Legs[3].IsInContact() {

		return true
	}

	return false
}

func (odometry *Odometry) noFeetInContact() bool {
	if !odometry.quadBase.Legs[0].IsInContact() &&
		!odometry.quadBase.Legs[1].IsInContact() &&
		!odometry.quadBase.Legs[2].IsInContact() &&
		!odometry.quadBase.Legs[3].IsInContact() {

		return true
	}

	return false
}

func (odometry *Odometry) GetVelocities(velocities cstructs.Velocities, currentTime time.Time) cstructs.Velocities {
	if odometry.allFeetInContact() || odometry.noFeetInContact() {
		velocities.Linear.SetX(0.0)
		velocities.Linear.SetY(0.0)
		velocities.Angular.SetZ(0.0)

		odometry.previousVelocity.Linear.SetX(0.0)
		odometry.previousVelocity.Linear.SetY(0.0)
		odometry.previousVelocity.Angular.SetZ(0.0)

		return velocities
	}

	totalContract := 0
	xSum := float32(0.0)
	ySum := float32(0.0)
	thetaSum := float32(0.0)

	for i := 0; i < 4; i++ {
		currentFootPosition := odometry.quadBase.Legs[i].FootFromBase()

		footInContact := odometry.quadBase.Legs[i].IsInContact()

		deltaX := odometry.previousFootPositions[i].X() - currentFootPosition.X()
		deltaY := odometry.previousFootPositions[i].Y() - currentFootPosition.Y()

		currentTheta := math32.Atan2(currentFootPosition.X(), currentFootPosition.Y())
		deltaTheta := currentTheta - odometry.previousTheta[i]

		if footInContact {
			totalContract = totalContract + 1
			thetaSum = thetaSum + deltaTheta
			xSum = xSum + (deltaX / 2)
			ySum = ySum + (deltaY / 2)
		}

		odometry.previousFootPositions[i] = currentFootPosition
		odometry.previousFootContacts[i] = footInContact
		odometry.previousTheta[i] = currentTheta
	}

	dt := float32(currentTime.Sub(odometry.previousTime).Microseconds())

	velocities.Linear.SetX(((1 - odometry.beta) * ((xSum * odometry.quadBase.GaitConfig().OdomScalar) / dt)) + (odometry.beta * odometry.previousVelocity.Linear.X()))
	velocities.Linear.SetY(((1 - odometry.beta) * ((ySum * odometry.quadBase.GaitConfig().OdomScalar) / dt)) + (odometry.beta * odometry.previousVelocity.Linear.Y()))
	velocities.Angular.SetZ(((1 - odometry.beta) * (thetaSum / dt)) + (odometry.beta * odometry.previousVelocity.Angular.Z()))

	odometry.previousVelocity = velocities
	odometry.previousTime = currentTime

	return velocities
}
