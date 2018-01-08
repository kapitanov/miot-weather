package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type weather struct {
	Now float32 `json:"now"`
}

const (
	defaultCity   = "moscow"
	sleepDuration = time.Minute * 60
)

var (
	city           string
	currentWeather *weather
	lock           sync.Mutex
)

type yandexWeatherCliOutput struct {
	City string  `json:"city"`
	Now  float32 `json:"term_now"`
}

func weatherInit() error {
	city = os.Getenv("WEATHER_CITY")
	if city == "" {
		city = defaultCity
	}

	err := weatherUpdate()
	if err != nil {
		return err
	}

	go func() {
		for {
			time.Sleep(sleepDuration)
			err := weatherUpdate()
			if err != nil {
				log.Fatalf("Unable to fetch data! %s\n", err)
			}
		}
	}()

	return nil
}

func weatherUpdate() error {
	w, err := weatherQuery()
	if err != nil {
		return err
	}
	lock.Lock()
	defer lock.Unlock()

	currentWeather = w
	log.Printf("Now in %s %0f deg\n", city, w.Now)
	mqttPublish()

	return nil
}

func weatherGet() *weather {
	lock.Lock()
	defer lock.Unlock()

	return currentWeather
}

func weatherQuery() (*weather, error) {
	cmd := exec.Command("yandex-weather-cli", "--json", city)
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var output yandexWeatherCliOutput
	err = json.Unmarshal(stdout, &output)
	if err != nil {
		return nil, err
	}

	var w weather
	w.Now = output.Now

	return &w, nil
}
