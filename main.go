package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	lat          = "43.71403201642408"
	lon          = "-70.3336353833058"
	appid        = "1f7d337e1b1fee64656d6a18ab11e55b"
	exclude      = "daily,hourly,minutely"
	compensation = 273.15 // Kelvin to Celcius
)
const (
	cycleTime = 1 * time.Second
	runTime   = 4 * cycleTime
	debug     = false
)

type Finish struct {
	stop chan bool
	wg   *sync.WaitGroup
}

func main() {
	var wg sync.WaitGroup
	stop := make(chan bool, 0)
	fan, _ := MakeFan()

	doneRunning := time.After(runTime) // how the program normally ends
	go Supervise(stop, &wg, fan)
	<-doneRunning
	close(stop)
	wg.Wait()
}

func Supervise(stop chan bool, wg *sync.WaitGroup, fan Fan) {
	signals := make(chan bool)
	wg.Add(1)
	defer wg.Done() // LIFO; run this last
	defer fan.Off()
	nextCycle := time.After(0 * time.Second)
	for {
		select {
		case <-signals:
			return
			fmt.Printf("user quit\n")
		case <-stop:
			return
			fmt.Printf("stahp\n")
		case <-nextCycle:
			nextCycle = time.After(cycleTime)
			cycle(stop, wg, fan)

			select {
			case <-nextCycle:
				panic(nil)
			default:
				continue
			}
		}
	}
}

func cycle(stop chan bool, wg *sync.WaitGroup, fan Fan) {
	airOut, err := GetOutsideAir()
	airIn, err := GetInsideAir()
	if err != nil {
		return
	}
	airOut.ControlFan(airIn, fan)
	fmt.Printf("Air Outside: %v, Air Inside: %v, Fan: %v\n", airOut.String(), airIn.String(), fan.String())
}
