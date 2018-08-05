package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// Config holds the configuration required for the Soil Monitor module.
type Config struct {
	LocationName string  `json:"locationName"` // Name of the Location
	LocationID   string  `json:"locationID"`   // Location identifier
	Latitude     float32 `json:"latitude"`     // Location Latitude
	Longitude    float32 `json:"longitude"`    // Location Longitude
	Provider     int     `json:"provider"`     // Provider type: 0=OpenWeather, 1=AccuWeather
	UnitType     int     `json:"unitType"`     // Unit type: 0=Metric, 1=Imperial
	AppID        string  `json:"appID"`        // Provider Application Identifier
}

// ReadFromFile will read the configuration settings from the specified file
func (c *Config) ReadFromFile(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		b, err := ioutil.ReadFile(path)
		if err == nil {
			err = json.Unmarshal(b, &c)
		}
	}
	c.SetDefaults()
	return err
}

// WriteToFile will write the configuration settings to the specified file
func (c *Config) WriteToFile(path string) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0666)
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (c *Config) ReadFrom(r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			err = json.Unmarshal(b, &c)
		}
	}
	c.SetDefaults()
	return err
}

// WriteTo serializes the entity and writes it to the http response
func (c *Config) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(c)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}

// Serialize serializes the entity and returns the serialized string
func (c *Config) Serialize() (string, error) {
	b, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Deserialize deserializes the specified string into the entity values
func (c *Config) Deserialize(v string) error {
	err := json.Unmarshal([]byte(v), &c)
	c.SetDefaults()
	return err
}

// SetDefaults checks the configuration and makes sure that, if a value is not configured, the default value is set.
func (c *Config) SetDefaults() {
	// Set any defaults required
	if c.Longitude == 0 && c.Latitude == 0 {
		i, err := GetIPLocationInfo()
		if err == nil {
			c.LocationName = fmt.Sprintf("%s, %s", i.City, i.Country)
			c.Longitude = i.Longitude
			c.Latitude = i.Latitude

		}
	}
}
