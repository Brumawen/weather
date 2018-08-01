package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// WeatherController handles the Web Methods for retrieving weather and forecast information.
type WeatherController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *WeatherController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
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
		// Check to see if we have already downloaded the latest forecast
		lf := Forecast{}
		if _, err := os.Stat("lastforecast.json"); err == nil {
			// File exists, check if this provider created it
			// and whether it was created less than 1 hour ago and, if so, return this record
			if err = lf.ReadFromFile("lastforecast.json"); err == nil {
				if lf.Current.Provider == p.GetProviderName() && time.Since(lf.Current.Created).Minutes() <= 60 {
					c.LogInfo("Returning cached forecast.")
					if err := lf.WriteTo(w); err != nil {
						c.LogError("Error serializing forecast information. " + err.Error())
					} else {
						return
					}
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
