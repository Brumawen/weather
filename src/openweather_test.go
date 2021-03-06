package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGetOpenWeather(t *testing.T) {
	c := Config{}
	err := c.ReadFromFile("openweather_test.json")
	if err != nil {
		t.Error(err)
	}

	o := OpenWeather{Config: &c}

	w, err := o.GetWeather()
	if err != nil {
		t.Error(err)
	}
	data, err := json.Marshal(&w)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("test.json", data, 0666)
	if err != nil {
		t.Error(err)
	}
}

func TestGetOpenWeatherForecast(t *testing.T) {
	c := Config{}
	err := c.ReadFromFile("openweather_test.json")
	if err != nil {
		t.Error(err)
	}

	o := OpenWeather{Config: &c}

	f, err := o.GetForecast()
	if err != nil {
		t.Error(err)
	}
	data, err := json.Marshal(&f)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("test.json", data, 0666)
	if err != nil {
		t.Error(err)
	}
}
