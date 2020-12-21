package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/examples/lib/dev"
)

const (
	lat          = "43.71403201642408"
	lon          = "-70.3336353833058"
	appid        = "1f7d337e1b1fee64656d6a18ab11e55b"
	exclude      = "daily,hourly,minutely"
	compensation = 273.15 // Kelvin to Celcius
)
const (
	cycleTime = 30 * time.Second
	runTime   = 4 * cycleTime
	debug     = false
)

type Finish struct {
	stop chan bool
	wg   *sync.WaitGroup
}

var device = "hci0"

func main() {
	var wg sync.WaitGroup
	stop := make(chan bool, 0)
	fan, _ := MakeFan()

	device, err := dev.NewDevice(device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(device)

	doneRunning := time.After(runTime) // how the program normally ends
	go Supervise(stop, &wg, fan, device)
	<-doneRunning
	close(stop)
	wg.Wait()
}

func Supervise(stop chan bool, wg *sync.WaitGroup, fan Fan, device ble.Device) {
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
			fmt.Printf("stahp\n")
			return
		case <-nextCycle:
			nextCycle = time.After(cycleTime)
			cycle(stop, wg, fan, device)

			select {
			case <-nextCycle:
				fmt.Printf("Cycle time shorter than cycle duration!\n")
			default:
				continue
			}
		}
	}
}

func cycle(stop chan bool, wg *sync.WaitGroup, fan Fan, device ble.Device) {
	airOut, err := GetOutsideAir()
	airIn, err := GetInsideAir(device)
	if err != nil {
		return
	}
	airOut.ControlFan(airIn, fan)
	fmt.Printf("Air Outside: %v, Air Inside: %v, Fan: %v\n", airOut.String(), airIn.String(), fan.String())
}
