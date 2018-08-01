package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGetAccuWeather(t *testing.T) {
	c := Config{}
	err := c.ReadFromFile("accuweather_test.json")
	if err != nil {
		t.Error(err)
	}

	o := AccuWeather{Config: &c}

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

func TestGetAccuWeatherForecast(t *testing.T) {
	c := Config{}
	err := c.ReadFromFile("accuweather_test.json")
	if err != nil {
		t.Error(err)
	}

	o := AccuWeather{Config: &c}

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
