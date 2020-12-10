package main

import "fmt"

const (
	roomTemperature float64 = 20.0
	roomHumidity    float64 = 40.0
)

type Air struct {
	Temperature float64
	Humidity    float64
}

func (a *Air) String() string {
	return fmt.Sprintf("%.1f C and %.0f%% (%.2fg/L)", a.Temperature, a.Humidity, a.WaterGramsPerLiter())
}

func NewAir(temp, humidity float64) Air {
	return Air{Temperature: temp, Humidity: humidity}
}

func MakeNewAir() Air {
	return Air{Temperature: roomTemperature, Humidity: roomHumidity}
}
