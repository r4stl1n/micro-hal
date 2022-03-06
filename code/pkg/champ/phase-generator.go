package champ

import (
	"github.com/r4stl1n/micro-hal/code/pkg/champ/cbase"
	"time"
)

type PhaseGenerator struct {
	time time.Time

	base              *cbase.QuadBase
	lastTouchDown     time.Time
	hasSwung          bool
	hasStarted        bool
	stancePhaseSignal [4]float32
	swingPhaseSignal  [4]float32
}

func (phaseGenerator *PhaseGenerator) Init(quadBase *cbase.QuadBase, currentTime time.Time) *PhaseGenerator {

	*phaseGenerator = PhaseGenerator{
		base:              quadBase,
		lastTouchDown:     currentTime,
		hasSwung:          false,
		hasStarted:        false,
		stancePhaseSignal: [4]float32{0.0, 0.0, 0.0, 0.0},
		swingPhaseSignal:  [4]float32{0.0, 0.0, 0.0, 0.0},
	}

	return phaseGenerator
}

func (phaseGenerator *PhaseGenerator) Run(targetVelocity float32, stepLength float32, currentTime time.Time) {

	secondsToMicro := float32(1000000)

	elapsedTimeRef := float32(0.0)
	swingPhasePeriod := 0.25 * secondsToMicro
	legClocks := [4]float32{0.0, 0.0, 0.0, 0.0}
	stancePhasePeriod := phaseGenerator.base.GaitConfig().StanceDuration * secondsToMicro
	stridePeriod := stancePhasePeriod + swingPhasePeriod

	if targetVelocity == 0.0 {
		elapsedTimeRef = 0
		phaseGenerator.lastTouchDown = time.Time{}
		phaseGenerator.hasSwung = false

		for i := 0; i < 4; i++ {
			legClocks[i] = 0.0
			phaseGenerator.stancePhaseSignal[i] = 0.0
			phaseGenerator.swingPhaseSignal[i] = 0.0
		}
		return
	}

	if !phaseGenerator.hasStarted {
		phaseGenerator.hasStarted = true
		phaseGenerator.lastTouchDown = currentTime
	}

	if float32(currentTime.Sub(phaseGenerator.lastTouchDown).Microseconds()) >= stridePeriod {
		phaseGenerator.lastTouchDown = currentTime
	}

	if elapsedTimeRef >= stridePeriod {
		elapsedTimeRef = stridePeriod
	} else {
		elapsedTimeRef = float32(currentTime.Sub(phaseGenerator.lastTouchDown).Microseconds())
	}

	legClocks[0] = elapsedTimeRef - (0.0 * stridePeriod)
	legClocks[1] = elapsedTimeRef - (0.5 * stridePeriod)
	legClocks[2] = elapsedTimeRef - (0.5 * stridePeriod)
	legClocks[3] = elapsedTimeRef - (0.5 * stridePeriod)

	for i := 0; i < 4; i++ {

		if legClocks[i] > 0 && legClocks[i] < stancePhasePeriod {
			phaseGenerator.stancePhaseSignal[i] = legClocks[i] / stancePhasePeriod
		} else {
			phaseGenerator.stancePhaseSignal[i] = 0
		}

		if legClocks[i] > -swingPhasePeriod && legClocks[i] < 0 {
			phaseGenerator.swingPhaseSignal[i] = (legClocks[i] + swingPhasePeriod) / swingPhasePeriod
		} else if legClocks[i] > swingPhasePeriod && legClocks[i] < stridePeriod {
			phaseGenerator.swingPhaseSignal[i] = (legClocks[i] - stancePhasePeriod) / swingPhasePeriod
		} else {
			phaseGenerator.swingPhaseSignal[i] = 0
		}
	}

	if !phaseGenerator.hasSwung && phaseGenerator.stancePhaseSignal[0] < 0.5 {
		phaseGenerator.stancePhaseSignal[0] = 0.0
		phaseGenerator.stancePhaseSignal[3] = 0.0
		phaseGenerator.stancePhaseSignal[1] = 0.0
		phaseGenerator.stancePhaseSignal[2] = 0.0
	} else {
		phaseGenerator.hasSwung = true
	}
}
