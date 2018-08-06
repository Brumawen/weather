package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// WeatherController handles the Web Methods for retrieving weather and forecast information.
type WeatherController struct {
	Srv *Server
}

// WeatherPageData holds the data used to populate the weather html page
type WeatherPageData struct {
	UnitType      int           // Unit Type: 0=Metric, 1=Imperial
	Temp          float32       // Current Temperature
	Pressure      float32       // Current Pressure
	Humidity      float32       // Current Humidity
	WindSpeed     float32       // Wind Speed
	WindDirection float32       // Wind Direction
	Sunrise       time.Time     // Time of sunrise
	Sunset        time.Time     // Time of sunset
	ReadTime      time.Time     // Time reading was taken
	WeatherIcon   string        // Weather Icon
	MoonIcon      string        // Moon Icon
	MoonDesc      string        // Moon Description
	Forecast      []ForecastDay // Forecast
}

// AddController adds the controller routes to the router
func (c *WeatherController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Path("/weather.html").Handler(http.HandlerFunc(c.handleWeatherWebPage))
	router.Path("/dashboard.html").Handler(http.HandlerFunc(c.handleWeatherWebPage))
	router.Methods("GET").Path("/weather/current").Name("GetCurrent").
		Handler(Logger(c, http.HandlerFunc(c.handleGetCurrent)))
	router.Methods("GET").Path("/weather/forecast").Name("GetForecast").
		Handler(Logger(c, http.HandlerFunc(c.handleGetForecast)))
}

// LogInfo is used to log information messages for this controller.
func (c *WeatherController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("WeatherController: ", a[1:len(a)-1])
}

// LogError is used to log error messages for this controller.
func (c *WeatherController) LogError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("WeatherController: ", a[1:len(a)-1])
}

func (c *WeatherController) handleWeatherWebPage(w http.ResponseWriter, r *http.Request) {
	p, err := c.getWeatherProvider()
	if err != nil {
		c.LogError("Error getting weather provider." + err.Error())
		http.Error(w, "Error getting weather prvider. "+err.Error(), 500)
		return
	}

	cf := c.getCurrentForecast(p)

	v := WeatherPageData{
		UnitType:      c.Srv.Config.UnitType,
		Temp:          cf.Current.Temp,
		Pressure:      cf.Current.Pressure,
		Humidity:      cf.Current.Humidity,
		WindSpeed:     cf.Current.WindSpeed,
		WindDirection: cf.Current.WindDirection,
		Sunrise:       cf.Current.Sunrise,
		Sunset:        cf.Current.Sunset,
		Forecast:      cf.Forecast,
	}
	v.MoonIcon, v.MoonDesc = c.getMoonIconInfo()

	t := template.Must(template.ParseFiles("./html/weather.html"))
	t.Execute(w, v)
}

// Get the current weather information
func (c *WeatherController) handleGetCurrent(w http.ResponseWriter, r *http.Request) {
	if p, err := c.getWeatherProvider(); err != nil {
		c.LogError("Error getting weather provider. " + err.Error())
		http.Error(w, "Error getting weather provider. "+err.Error(), 500)
	} else {
		// Check to see if we have already downloaded the latest weather
		lw := Weather{}
		if _, err := os.Stat("lastweather.json"); err == nil {
			// File exists, check if this provider created it
			// and whether it was created less than 1 hour ago and, if so, return this record
			if err = lw.ReadFromFile("lastweather.json"); err == nil {
				if lw.Provider == p.GetProviderName() && time.Since(lw.Created).Minutes() <= 60 {
					c.LogInfo("Returning cached weather.")
					if err := lw.WriteTo(w); err != nil {
						c.LogError("Error serializing weather information. " + err.Error())
					} else {
						return
					}
				}
			}
		}

		// Get the weather information from the weather site
		cw, err := p.GetWeather()
		if err != nil {
			c.LogError("Error getting weather information. " + err.Error())
			cw = lw
		} else {
			cw.WriteToFile("lastweather.json")
		}
		if err := cw.WriteTo(w); err != nil {
			c.LogError("Error serializing weather information. " + err.Error())
			http.Error(w, "Error serializing weather information. "+err.Error(), 500)
		}
	}
}

// Get the current forecast information
func (c *WeatherController) handleGetForecast(w http.ResponseWriter, r *http.Request) {
	if p, err := c.getWeatherProvider(); err != nil {
		c.LogError("Error getting weather provider. " + err.Error())
		http.Error(w, "Error getting weather provider. "+err.Error(), 500)
	} else {
		cf := c.getCurrentForecast(p)
		if err := cf.WriteTo(w); err != nil {
			c.LogError("Error serializing forecast information. " + err.Error())
			http.Error(w, "Error serializing forecast information. "+err.Error(), 500)
		}
	}
}

func (c *WeatherController) getWeatherProvider() (WeatherProvider, error) {
	switch c.Srv.Config.Provider {
	case 0:
		// OpenWeather
		ow := new(OpenWeather)
		ow.SetConfig(c.Srv.Config)
		return ow, nil
	case 1:
		// Accuweather
		aw := new(AccuWeather)
		aw.SetConfig(c.Srv.Config)
		return aw, nil
	default:
		return nil, errors.New("Invalid Weather provider")
	}
}

func (c *WeatherController) getCurrentForecast(p WeatherProvider) Forecast {
	// Check to see if we have already downloaded the latest forecast
	lf := Forecast{}
	if _, err := os.Stat("lastforecast.json"); err == nil {
		// File exists, check if this provider created it
		// and whether it was created less than 1 hour ago and, if so, return this record
		if err = lf.ReadFromFile("lastforecast.json"); err == nil {
			if lf.Current.Provider == p.GetProviderName() && time.Since(lf.Current.Created).Minutes() <= 60 {
				c.LogInfo("Returning cached forecast.")
				return lf
			}
		}
	}

	c.LogInfo("Getting fresh forecast from the provider.")
	cf, err := p.GetForecast()
	if err != nil {
		c.LogError("Error getting forecast information. " + err.Error())
		cf = lf
	} else {
		cf.WriteToFile("lastforecast.json")
	}
	return cf
}

func (c *WeatherController) getWeatherIconInfo(i int, day bool) string {
	// Icon numbers:
	// 1 = Clear sky
	// 2 = Scattered clouds
	// 3 = Partly cloudy
	// 4 = Cloudy
	// 5 = Scattered Rain
	// 6 = Rain
	// 7 = Thunderstorms
	// 8 = Snow
	// 9 = Mist/ Fog
	if day {
		switch i {
		case 1:
			return "wi-day-sunny"
		case 2:
			return "wi-day-cloudy"
		case 3:
			return "wi-cloud"
		case 4:
			return "wi-cloudy"
		case 5:
			return "wi-day-rain"
		case 6:
			return "wi-rain"
		case 7:
			return "wi-thunderstorm"
		case 8:
			return "wi-snow"
		case 9:
			return "wi-dust"
		}
	} else {
		switch i {
		case 1:
			return "wi-night-clear"
		case 2:
			return "wi-night-alt-cloudy"
		case 3:
			return "wi-cloud"
		case 4:
			return "wi-cloudy"
		case 5:
			return "wi-night-alt-rain"
		case 6:
			return "wi-rain"
		case 7:
			return "wi-thunderstorm"
		case 8:
			return "wi-snow"
		case 9:
			return "wi-dust"
		}
	}

	// Return this if we have an icon we don't know how to deal with
	return "wi-alien"
}

func (c *WeatherController) getMoonIconInfo() (string, string) {
	m := Moon{}
	m.ForDate(time.Now())
	switch int(m.Age) {
	case 0:
		return "wi-moon-alt-new", "New"
	case 1:
		return "wi-moon-alt-waxing-crescent-1", "Waxing Crescent"
	case 2:
		return "wi-moon-alt-waxing-crescent-2", "Waxing Crescent"
	case 3:
		return "wi-moon-alt-waxing-crescent-3", "Waxing Crescent"
	case 4:
		return "wi-moon-alt-waxing-crescent-4", "Waxing Crescent"
	case 5:
		return "wi-moon-alt-waxing-crescent-5", "Waxing Crescent"
	case 6:
		return "wi-moon-alt-waxing-crescent-6", "Waxing Crescent"
	case 7:
		return "wi-moon-alt-first-quarter", "First Quarter"
	case 8:
		return "wi-moon-alt-waxing-gibbous-1", "Waxing Gibbous"
	case 9:
		return "wi-moon-alt-waxing-gibbous-2", "Waxing Gibbous"
	case 10:
		return "wi-moon-alt-waxing-gibbous-3", "Waxing Gibbous"
	case 11:
		return "wi-moon-alt-waxing-gibbous-4", "Waxing Gibbous"
	case 12:
		return "wi-moon-alt-waxing-gibbous-5", "Waxing Gibbous"
	case 13:
		return "wi-moon-alt-waxing-gibbous-6", "Waxing Gibbous"
	case 14:
		return "wi-moon-alt-full", "Full"
	case 15:
		return "wi-moon-alt-waning-gibbous-1", "Waning Gibbous"
	case 16:
		return "wi-moon-alt-waning-gibbous-2", "Waning Gibbous"
	case 17:
		return "wi-moon-alt-waning-gibbous-3", "Waning Gibbous"
	case 18:
		return "wi-moon-alt-waning-gibbous-4", "Waning Gibbous"
	case 19:
		return "wi-moon-alt-waning-gibbous-5", "Waning Gibbous"
	case 20:
		return "wi-moon-alt-waning-gibbous-6", "Waning Gibbous"
	case 21:
		return "wi-moon-alt-third-quarter", "Third Quarter"
	case 22:
		return "wi-moon-alt-waning-crescent-1", "Waning Crescent"
	case 23:
		return "wi-moon-alt-waning-crescent-2", "Waning Crescent"
	case 24:
		return "wi-moon-alt-waning-crescent-3", "Waning Crescent"
	case 25:
		return "wi-moon-alt-waning-crescent-4", "Waning Crescent"
	case 26:
		return "wi-moon-alt-waning-crescent-5", "Waning Crescent"
	default:
		return "wi-moon-alt-waning-crescent-6", "Waning Crescent"
	}
}
