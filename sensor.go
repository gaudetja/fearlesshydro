package main

import (
	"fmt"
	"log"
	"time"

	"github.com/go-ble/ble"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

const (
	moistTemp     = 23.0
	moistHumidity = 77.0

	fineTemp     = 18
	fineHumidity = 40.0
)

var (
	addr = "a4:c1:38:ab:4d:cc"
	dup  = true
	du   = time.Second
)

func advHandler(a ble.Advertisement) {
	return
	if a.Connectable() {
		fmt.Printf("[%s] C %3d:", a.Addr(), a.RSSI())
	} else {
		fmt.Printf("[%s] N %3d:", a.Addr(), a.RSSI())
	}
	comma := ""
	if len(a.LocalName()) > 0 {
		fmt.Printf(" Name: %s", a.LocalName())
		comma = ","
	}
	if len(a.Services()) > 0 {
		fmt.Printf("%s Svcs: %v", comma, a.Services())
		comma = ","
	}
	if len(a.ManufacturerData()) > 0 {
		fmt.Printf("%s MD: %X", comma, a.ManufacturerData())
	}
	fmt.Printf("\n")
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}

func propString(p ble.Property) string {
	var s string
	for k, v := range map[ble.Property]string{
		ble.CharBroadcast:   "B",
		ble.CharRead:        "R",
		ble.CharWriteNR:     "w",
		ble.CharWrite:       "W",
		ble.CharNotify:      "N",
		ble.CharIndicate:    "I",
		ble.CharSignedWrite: "S",
		ble.CharExtended:    "E",
	} {
		if p&k != 0 {
			s += v
		}
	}
	return s
}

type airGetter struct {
	out chan Air
}

func (ag airGetter) notifyHandler(req []byte) {
	temp := (float64(req[1])*256 + float64(req[0])) / 100
	humidity := float64(req[2])
	voltage := (float64(req[4])*256 + float64(req[3])) / 1000
	// fmt.Printf("Got message: %v\n", req)
	fmt.Printf("temp: %vC, hum: %v%%, batt: %vV\n", temp, humidity, voltage)
	ag.out <- Air{Temperature: temp, Humidity: humidity}
}

func GetInsideAir(dev ble.Device) (Air, error) {

	filter := func(a ble.Advertisement) bool {
		return a.Addr().String() == addr
	}
	cln, err := ble.Connect(context.Background(), filter)
	if err != nil {
		log.Fatalf("can't connect : %s", err)
	}
	serviceID, _ := ble.Parse("ebe0ccb07a0a4b0c8a1a6ff2997da3a6")
	characteristicID, _ := ble.Parse("ebe0ccc17a0a4b0c8a1a6ff2997da3a6")
	var air = Air{Temperature: moistTemp, Humidity: moistHumidity}

	p, err := cln.DiscoverProfile(true)
	if err != nil {
		log.Fatalf("can't discover profile: %s", err)
	}
	for _, s := range p.Services {
		if !s.UUID.Equal(serviceID) {
			continue
		}
		for _, c := range s.Characteristics {
			if !c.UUID.Equal(characteristicID) {
				continue
			}
			if (c.Property & ble.CharRead) != 0 {
				// b, err := cln.ReadCharacteristic(c)
				// if err != nil {
				// 	fmt.Printf("Failed to read characteristic: %s\n", err)
				// 	continue
				// }
				airChan := make(chan Air)
				getter := airGetter{out: airChan}
				err = cln.Subscribe(c, false, getter.notifyHandler)
				if err != nil {
					fmt.Printf("While sub: %v\n", err)
				}
				timeout := time.After(10 * time.Second)
				select {
				case <-timeout:
				case air = <-airChan:
				}
				err = cln.Unsubscribe(c, false)
				if err != nil {
					fmt.Printf("While unsub: %v\n", err)
				}
			}
		}
	}

	cln.CancelConnection()

	return air, nil
}
