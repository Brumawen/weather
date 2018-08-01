package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Weather holds the current weather information
type Weather struct {
	Provider      string    `json:"provider"`      // Provider
	Created       time.Time `json:"created"`       // Date and time the information was created by the provider
	ID            string    `json:"locationID"`    // Location ID
	Name          string    `json:"locationName"`  // Location Name
	Temp          float32   `json:"temp"`          // Current Temperature
	Pressure      float32   `json:"pressure"`      // Current Pressure
	Humidity      float32   `json:"humidity"`      // Current Humidity
	WindSpeed     float32   `json:"windSpeed"`     // Current Wind Speed
	WindDirection float32   `json:"windDirection"` // Current Wind Direction
	Icon          string    `json:"icon"`          // Weather Icon
	ReadingTime   time.Time `json:"readingTime"`   // Date and Time the reading was taken
	Sunrise       time.Time `json:"sunrise"`       // Time of Sunrise
	Sunset        time.Time `json:"sunset"`        // Time of Sunset
}

// Forecast holds the current weather and the forecast weather information
type Forecast struct {
	Current Weather       // Current Weather
	Days    []ForecastDay // Weather Forecast
}

// ForecastDay holds the temperature and weather forecase for a particular day
type ForecastDay struct {
	Day     time.Time // Forecast date
	Name    string    // Name of the day
	TempMin float32   // Minimum Temperature
	TempMax float32   // Maximum Temperature
	Icon    string    // Weather Icon
}

// ReadFromFile will read the weather information from the specified file
func (c *Weather) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &c)
		}
	}
	return err
}

// WriteToFile will write the weather information to the specified file
func (c *Weather) WriteToFile(path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

// WriteTo serializes the entity and writes it to the http response
func (c *Weather) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// ReadFromFile will read the forecast information from the specified file
func (c *Forecast) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &c)
		}
	}
	return err
}

// WriteToFile will write the forecast information to the specified file
func (c *Forecast) WriteToFile(path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

// WriteTo serializes the entity and writes it to the http response
func (c *Forecast) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}
