package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GetOutsideAir() (Air, error) {
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/onecall?lat=%s&lon=%s&exclude=%s&appid=%s", lat, lon, exclude, appid)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Bad resp: %v, %v", resp, err)
	}
	defer resp.Body.Close()
	rjson := make(map[string]interface{})

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Bad body: %v\n\n %v", resp, err)
		return MakeNewAir(), err
	}
	err = json.Unmarshal(body, &rjson)
	if debug {
		out := bufio.NewWriter(os.Stdout)
		enc := json.NewEncoder(out)
		enc.SetIndent("", "    ")
		if err := enc.Encode(rjson); err != nil {
			fmt.Println("Error decoding output: %v\n\n %v", resp, err)
			return MakeNewAir(), err
		}
		out.Flush()
	}

	current := rjson["current"].(map[string]interface{})

	airNow := Air{
		Humidity:    current["humidity"].(float64),
		Temperature: current["temp"].(float64) - compensation,
	}
	return airNow, nil
}
