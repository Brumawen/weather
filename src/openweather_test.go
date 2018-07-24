package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestGetWeather(t *testing.T) {
	o := OpenWeather{
		AppID:     "4e2ff3184d5ce080b27167bb5544beb0",
		Latitude:  -25.989,
		Longitude: 28.003,
	}

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

func TestGetForecast(t *testing.T) {
	o := OpenWeather{
		AppID:     "4e2ff3184d5ce080b27167bb5544beb0",
		Latitude:  -25.989,
		Longitude: 28.003,
	}

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
