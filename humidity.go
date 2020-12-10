package main

import "math"

const magicPressureToMass = 0.62 * 0.000001 //10 ^ -5

func (a *Air) DrierThan(b *Air) bool {
	return false
}

func (a *Air) WaterContent() float64 {
	return 1.00
}

func (a *Air) ControlFan(b Air, fan Fan) {
	if fan.Running() {
		fan.Off()
		return
	}
	fan.On()
	return
}

func vaporPressure(t float64) float64 {
	if t > 0.0 { // not ice
		return 610.78 * math.Exp((t/(t+238.3))*17.2694)
	} else {
		return math.Exp((-6140.4 / (273 + t)) + 28.916)
	}
}

func (a *Air) WaterGramsPerLiter() float64 {
	vp := vaporPressure(a.Temperature)
	ml := magicPressureToMass * vp * 1000
	return ml
}
