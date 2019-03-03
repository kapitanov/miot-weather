package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
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
	cityName         string
	currentWeather   *weather
	lock             sync.Mutex
	temperatureRegex *regexp.Regexp
)

type yandexWeatherCliOutput struct {
	City string  `json:"city"`
	Now  float32 `json:"term_now"`
}

func weatherInit() error {
	cityName = os.Getenv("WEATHER_CITY")
	if cityName == "" {
		cityName = defaultCity
	}

	temperatureRegex = regexp.MustCompile(`(?i)<span class="temp__value">(.*?)</span>`)

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
	log.Printf("Now in %s %0f deg\n", cityName, w.Now)
	mqttPublish()

	return nil
}

func weatherGet() *weather {
	lock.Lock()
	defer lock.Unlock()

	return currentWeather
}

func weatherQuery() (*weather, error) {
	response, err := http.Get("https://yandex.ru/pogoda/" + cityName)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	stringContent := string(content)
	str := temperatureRegex.FindStringSubmatch(stringContent)[1]
	str = strings.Replace(str, "\u2212", "-", 1)
	temperature, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return nil, err
	}

	var w weather
	w.Now = float32(temperature)

	return &w, nil
}
