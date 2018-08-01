package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// AccuWeather is an interface to the AccuWeatherMap internet API
type AccuWeather struct {
	Config *Config // Current Configuration
}

type accuWeatherResponse []struct {
	LocalObservationDateTime time.Time `json:"LocalObservationDateTime"`
	EpochTime                int       `json:"EpochTime"`
	WeatherText              string    `json:"WeatherText"`
	WeatherIcon              int       `json:"WeatherIcon"`
	IsDayTime                bool      `json:"IsDayTime"`
	Temperature              struct {
		Metric struct {
			Value    float64 `json:"Value"`
			Unit     string  `json:"Unit"`
			UnitType int     `json:"UnitType"`
		} `json:"Metric"`
		Imperial struct {
			Value    float64 `json:"Value"`
			Unit     string  `json:"Unit"`
			UnitType int     `json:"UnitType"`
		} `json:"Imperial"`
	} `json:"Temperature"`
	RelativeHumidity float64 `json:"RelativeHumidity"`
	Wind             struct {
		Direction struct {
			Degrees   float64 `json:"Degrees"`
			Localized string  `json:"Localized"`
			English   string  `json:"English"`
		} `json:"Direction"`
		Speed struct {
			Metric struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Metric"`
			Imperial struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Imperial"`
		} `json:"Speed"`
	} `json:"Wind"`
	WindGust struct {
		Speed struct {
			Metric struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Metric"`
			Imperial struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Imperial"`
		} `json:"Speed"`
	} `json:"WindGust"`
	UVIndex     float64 `json:"UVIndex"`
	UVIndexText string  `json:"UVIndexText"`
	CloudCover  float64 `json:"CloudCover"`
	Pressure    struct {
		Metric struct {
			Value    float64 `json:"Value"`
			Unit     string  `json:"Unit"`
			UnitType int     `json:"UnitType"`
		} `json:"Metric"`
		Imperial struct {
			Value    float64 `json:"Value"`
			Unit     string  `json:"Unit"`
			UnitType int     `json:"UnitType"`
		} `json:"Imperial"`
	} `json:"Pressure"`
	PressureTendency struct {
		LocalizedText string `json:"LocalizedText"`
		Code          string `json:"Code"`
	} `json:"PressureTendency"`
}

type accuForecastResponse struct {
	Headline struct {
		EffectiveDate      time.Time   `json:"EffectiveDate"`
		EffectiveEpochDate int         `json:"EffectiveEpochDate"`
		Severity           int         `json:"Severity"`
		Text               string      `json:"Text"`
		Category           string      `json:"Category"`
		EndDate            interface{} `json:"EndDate"`
		EndEpochDate       interface{} `json:"EndEpochDate"`
		MobileLink         string      `json:"MobileLink"`
		Link               string      `json:"Link"`
	} `json:"Headline"`
	DailyForecasts []struct {
		Date        time.Time `json:"Date"`
		EpochDate   int       `json:"EpochDate"`
		Temperature struct {
			Minimum struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Minimum"`
			Maximum struct {
				Value    float64 `json:"Value"`
				Unit     string  `json:"Unit"`
				UnitType int     `json:"UnitType"`
			} `json:"Maximum"`
		} `json:"Temperature"`
		Day struct {
			Icon       int    `json:"Icon"`
			IconPhrase string `json:"IconPhrase"`
		} `json:"Day"`
		Night struct {
			Icon       int    `json:"Icon"`
			IconPhrase string `json:"IconPhrase"`
		} `json:"Night"`
	} `json:"DailyForecasts"`
}

type accuLocationResponse struct {
	Version int    `json:"Version"`
	Key     string `json:"Key"`
	Code    string `json:"Code"`
	Message string `json:"Message"`
}

// SetConfig sets the configuration for the provider
func (p *AccuWeather) SetConfig(c *Config) {
	p.Config = c
}

// GetProviderName returns the name of the provider
func (p *AccuWeather) GetProviderName() string {
	return "AccuWeather"
}

// GetWeather returns the current weather for the location
func (p *AccuWeather) GetWeather() (Weather, error) {
	// Get the latest weather information from AccuWeather
	w := Weather{
		Provider: p.GetProviderName(),
		Created:  time.Now(),
	}
	if err := p.checkConfig(); err != nil {
		return w, err
	}
	w.ID = p.Config.LocationID

	url := fmt.Sprintf("http://dataservice.accuweather.com/currentconditions/v1/%s?apikey=%s&details=true", p.Config.LocationID, p.Config.AppID)
	resp, err := http.Get(url)
	if err == nil {
		err = p.decodeWeather(&w, resp.Body)
	}
	if err == nil {
		w.WriteToFile("lastweather.json")
	}
	return w, err
}

// GetForecast returns the current forecast for the location
func (p *AccuWeather) GetForecast() (Forecast, error) {
	// Get the latest weather information from AccuWeather
	f := Forecast{
		Current: Weather{
			Provider: p.GetProviderName(),
			Created:  time.Now(),
		},
	}
	if err := p.checkConfig(); err != nil {
		return f, err
	}

	metric := "true"
	if p.Config.Imperial {
		metric = "false"
	}
	url := fmt.Sprintf("http://dataservice.accuweather.com/forecasts/v1/daily/5day/%s?metric=%s&apikey=%s", p.Config.LocationID, metric, p.Config.AppID)
	resp, err := http.Get(url)
	if err == nil {
		err = p.decodeForecast(&f, resp.Body)
	}
	return f, err
}

func (p *AccuWeather) decodeWeather(w *Weather, r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			var r = accuWeatherResponse{}
			err = json.Unmarshal(b, &r)
			if err == nil && r != nil && len(r) != 0 {
				r1 := r[0]
				w.ReadingTime = r1.LocalObservationDateTime
				w.Humidity = float32(r1.RelativeHumidity)
				w.WindDirection = float32(r1.Wind.Direction.Degrees)
				if p.Config.Imperial {
					// Get imperial values
					w.Temp = float32(r1.Temperature.Imperial.Value)
					w.Pressure = float32(r1.Pressure.Imperial.Value)
					w.WindSpeed = float32(r1.Wind.Speed.Imperial.Value)
				} else {
					// Get metric values
					w.Temp = float32(r1.Temperature.Metric.Value)
					w.Pressure = float32(r1.Pressure.Metric.Value)
					w.WindSpeed = float32(r1.Wind.Speed.Metric.Value)
				}
			}
		}
	}
	return err
}

func (p *AccuWeather) decodeForecast(f *Forecast, r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			var r = accuForecastResponse{}
			err = json.Unmarshal(b, &r)
			if err == nil && r.DailyForecasts != nil && len(r.DailyForecasts) != 0 {
				for _, d := range r.DailyForecasts {
					fd := ForecastDay{}
					fd.Day = d.Date
					fd.Name = d.Date.Weekday().String()[:3]
					fd.TempMax = float32(d.Temperature.Maximum.Value)
					fd.TempMin = float32(d.Temperature.Minimum.Value)
					f.Days = append(f.Days, fd)
				}
			}
		}
	}
	return err
}

func (p *AccuWeather) checkConfig() error {
	if p.Config.AppID == "" {
		return errors.New("Accuweather API Key has not been set in the configuration")
	}
	if p.Config.LocationID == "" {
		url := fmt.Sprintf("http://dataservice.accuweather.com/locations/v1/cities/geoposition/search?apikey=%s&q=%f%%2C%f", p.Config.AppID, p.Config.Latitude, p.Config.Longitude)
		resp, err := http.Get(url)
		if err != nil {
			return errors.New("Error getting AccuWeather location information. " + err.Error())
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.New("Error reading AccuWeather location information. " + err.Error())
		}
		var r = accuLocationResponse{}
		err = json.Unmarshal(b, &r)
		if err != nil {
			return errors.New("Error deserializing AccuWeather location information. " + err.Error())
		}
		if r.Message != "" {
			return errors.New(r.Message)
		}
		p.Config.LocationID = r.Key
	}
	return nil
}