package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// OpenWeather is an interface to the OpenWeatherMap internet API
type OpenWeather struct {
	Config *Config // Current Configuration
}

type owWeatherResponse struct {
	Coord struct {
		Lon float32 `json:"lon"`
		Lat float32 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float32 `json:"temp"`
		Pressure float32 `json:"pressure"`
		Humidity float32 `json:"humidity"`
		TempMin  float32 `json:"temp_min"`
		TempMax  float32 `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float32 `json:"speed"`
		Deg   float32 `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float32 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}

type owForecastResponse struct {
	Cod     string  `json:"cod"`
	Message float32 `json:"message"`
	Cnt     int     `json:"cnt"`
	List    []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float32 `json:"temp"`
			TempMin   float32 `json:"temp_min"`
			TempMax   float32 `json:"temp_max"`
			Pressure  float32 `json:"pressure"`
			SeaLevel  float32 `json:"sea_level"`
			GrndLevel float32 `json:"grnd_level"`
			Humidity  float32 `json:"humidity"`
			TempKf    float32 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed float32 `json:"speed"`
			Deg   float32 `json:"deg"`
		} `json:"wind"`
		Sys struct {
			Pod string `json:"pod"`
		} `json:"sys"`
		DtTxt string `json:"dt_txt"`
	} `json:"list"`
	City struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Coord struct {
			Lat float32 `json:"lat"`
			Lon float32 `json:"lon"`
		} `json:"coord"`
		Country string `json:"country"`
	} `json:"city"`
}

// SetConfig sets the configuration for the provider
func (o *OpenWeather) SetConfig(c *Config) {
	o.Config = c
}

// GetProviderName returns the name of the provider
func (o *OpenWeather) GetProviderName() string {
	return "OpenWeather"
}

// GetWeather returns the current weather for the location
func (o *OpenWeather) GetWeather() (Weather, error) {
	w := Weather{
		Provider: o.GetProviderName(),
		Created:  time.Now(),
	}

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", o.Config.Latitude, o.Config.Longitude, o.Config.AppID)
	var resp, err = http.Get(url)
	if err == nil {
		err = o.decodeWeather(&w, resp.Body)
	}
	return w, err
}

// GetForecast returns the current forecast for the location
func (o *OpenWeather) GetForecast() (Forecast, error) {
	f := Forecast{
		Current: Weather{
			Provider: o.GetProviderName(),
			Created:  time.Now(),
		},
	}

	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/forecast?lat=%f&lon=%f&appid=%s&units=metric", o.Config.Latitude, o.Config.Longitude, o.Config.AppID)
	var resp, err = http.Get(url)
	if err == nil {
		err = o.decodeForecast(&f, resp.Body)
	}
	return f, err
}

// ReadFrom reads the string from the reader and deserializes it into the entity values
func (o *OpenWeather) decodeWeather(w *Weather, r io.ReadCloser) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			var resp = owWeatherResponse{}
			err = json.Unmarshal(b, &resp)
			if err == nil {
				w.ID = strconv.Itoa(resp.ID)
				w.Name = resp.Name
				if len(resp.Weather) != 0 {
					w.Icon = resp.Weather[0].Icon
				}
				w.Temp = resp.Main.Temp
				w.Humidity = resp.Main.Humidity
				w.Pressure = resp.Main.Pressure
				w.ReadingTime = time.Unix(int64(resp.Dt), 0)
				w.Sunrise = time.Unix(int64(resp.Sys.Sunrise), 0)
				w.Sunset = time.Unix(int64(resp.Sys.Sunset), 0)
				w.WindSpeed = resp.Wind.Speed
				w.WindDirection = resp.Wind.Deg
			}
		}
	}
	return err
}

func (o *OpenWeather) decodeForecast(f *Forecast, r io.Reader) error {
	b, err := ioutil.ReadAll(r)
	if err == nil {
		if b != nil && len(b) != 0 {
			var resp = owForecastResponse{}
			err = json.Unmarshal(b, &resp)
			if err == nil {
				if len(resp.List) != 0 {
					// Current weather
					cw := resp.List[0]
					f.Current = Weather{}
					f.Current.ID = string(resp.City.ID)
					f.Current.Name = resp.City.Name
					if len(cw.Weather) != 0 {
						f.Current.Icon = cw.Weather[0].Icon
					}
					f.Current.Temp = cw.Main.Temp
					f.Current.Humidity = cw.Main.Humidity
					f.Current.Pressure = cw.Main.Pressure
					ct := time.Unix(int64(cw.Dt), 0)
					f.Current.ReadingTime = ct
					f.Current.WindSpeed = cw.Wind.Speed
					f.Current.WindDirection = cw.Wind.Deg

					// Forecast
					cf := ForecastDay{}
					cf.Day = time.Date(ct.Year(), ct.Month(), ct.Day(), 0, 0, 0, 0, ct.Location())
					for _, i := range resp.List {
						ct = time.Unix(int64(i.Dt), 0)
						iDay := time.Date(ct.Year(), ct.Month(), ct.Day(), 0, 0, 0, 0, ct.Location())
						if iDay.Year() != cf.Day.Year() || iDay.YearDay() != cf.Day.YearDay() {
							// Date has changed
							f.Days = append(f.Days, cf)
							cf = ForecastDay{}
							cf.Day = time.Date(ct.Year(), ct.Month(), ct.Day(), 0, 0, 0, 0, ct.Location())
						}
						if cf.Icon == "" {
							cf.TempMin = i.Main.Temp
							cf.TempMax = i.Main.Temp
							cf.Day = ct
							cf.Name = ct.Weekday().String()[:3]
							if len(i.Weather) != 0 {
								cf.Icon = i.Weather[0].Icon
							} else {
								cf.Icon = "00d"
							}
						} else {
							if i.Main.Temp < cf.TempMin {
								cf.TempMin = i.Main.Temp
							}
							if i.Main.Temp > cf.TempMax {
								cf.TempMax = i.Main.Temp
							}
							if len(i.Weather) != 0 {
								icon := i.Weather[0].Icon
								a, err := strconv.Atoi(cf.Icon[:2])
								if err == nil {
									b, err := strconv.Atoi(icon[:2])
									if err == nil {
										if b > a {
											cf.Icon = icon
										}
									}
								}
							}
						}

					}
					f.Days = append(f.Days, cf)
				}
			}
		}
	}
	return err
}
