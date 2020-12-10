package main

const (
	moistTemp     = 23.0
	moistHumidity = 77.0

	fineTemp     = 18
	fineHumidity = 40.0
)

func GetInsideAir() (Air, error) {
	return Air{Temperature: moistTemp, Humidity: moistHumidity}, nil
}
