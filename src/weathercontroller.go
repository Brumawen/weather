package main

import (
	"fmt"
	"net/http"

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

// Get the current weather information
func (c *WeatherController) handleGetCurrent(w http.ResponseWriter, r *http.Request) {
	p := OpenWeather{
		Latitude:  c.Srv.Config.Latitude,
		Longitude: c.Srv.Config.Longitude,
		AppID:     c.Srv.Config.AppID,
	}
	i, err := p.GetWeather()
	if err != nil {
		http.Error(w, "Error getting weather information. "+err.Error(), 500)
	}
	if err := i.WriteTo(w); err != nil {
		http.Error(w, "Error serializing weather information. "+err.Error(), 500)
	}
}

// Get the current forecast information
func (c *WeatherController) handleGetForecast(w http.ResponseWriter, r *http.Request) {
	p := OpenWeather{
		Latitude:  c.Srv.Config.Latitude,
		Longitude: c.Srv.Config.Longitude,
		AppID:     c.Srv.Config.AppID,
	}
	i, err := p.GetForecast()
	if err != nil {
		http.Error(w, "Error getting forecast information. "+err.Error(), 500)
	}
	if err := i.WriteTo(w); err != nil {
		http.Error(w, "Error serializing forecast information. "+err.Error(), 500)
	}
}
