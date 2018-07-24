package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// Weather holds the current weather information
type Weather struct {
	ID            string    // Location ID
	Name          string    // Location Name
	Temp          float32   // Current Temperature
	Pressure      float32   // Current Pressure
	Humidity      float32   // Current Humidity
	WindSpeed     float32   // Current Wind Speed
	WindDirection float32   // Current Wind Direction
	Icon          string    // Weather Icon
	ReadingTime   time.Time // Date and Time the reading was taken
	Sunrise       time.Time // Time of Sunrise
	Sunset        time.Time // Time of Sunset
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
