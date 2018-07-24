package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// ConfigController handles the Web Methods for configuring the module.
type ConfigController struct {
	Srv *Server
}

// ConfigPageData holds the data used to write to the configuration page.
type ConfigPageData struct {
	LocationName string
	Longitude    string
	Latitude     string
	Provider     int
	AppID        string
}

// AddController adds the controller routes to the router
func (c *ConfigController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Path("/config.html").Handler(http.HandlerFunc(c.handleConfigWebPage))
	router.Methods("GET").Path("/config/get").Name("GetConfig").
		Handler(Logger(c, http.HandlerFunc(c.handleGetConfig)))
	router.Methods("POST").Path("/config/set").Name("SetConfig").
		Handler(Logger(c, http.HandlerFunc(c.handleSetConfig)))
}

func (c *ConfigController) handleConfigWebPage(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./html/config.html"))

	v := ConfigPageData{
		LocationName: c.Srv.Config.LocationName,
		Longitude:    fmt.Sprintf("%f", c.Srv.Config.Longitude),
		Latitude:     fmt.Sprintf("%f", c.Srv.Config.Latitude),
		Provider:     c.Srv.Config.Provider,
		AppID:        c.Srv.Config.AppID,
	}

	t.Execute(w, v)
}

func (c *ConfigController) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	if err := c.Srv.Config.WriteTo(w); err != nil {
		http.Error(w, "Error serializing configuration. "+err.Error(), 500)
	}
}

func (c *ConfigController) handleSetConfig(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	z, err := json.Marshal(r.Form)
	fmt.Println(string(z))

	nam := r.Form.Get("locationname")
	lon := r.Form.Get("longitude")
	lat := r.Form.Get("latitude")
	prv := r.Form.Get("provider")
	app := r.Form.Get("appid")

	if lon == "" {
		http.Error(w, "The Longitude of the forecast location must be specified", 500)
		return
	}
	a, err := strconv.ParseFloat(lon, 32)
	if err != nil {
		http.Error(w, "Invalid Longitude value", 500)
		return
	}
	if lat == "" {
		http.Error(w, "The Latitude of the forecast location must be specified", 500)
		return
	}
	b, err := strconv.ParseFloat(lat, 32)
	if err != nil {
		http.Error(w, "Invalid Latitude value", 500)
		return
	}
	if prv == "" {
		http.Error(w, "The Forecast Provider must be selected", 500)
		return
	}
	p, err := strconv.Atoi(prv)
	if err != nil || p != 0 {
		http.Error(w, "Invalid Forecast Provider value", 500)
		return
	}
	if app == "" {
		http.Error(w, "The Forecast Provider Application ID must be selected", 500)
		return
	}

	c.LogInfo("Setting new configuration values.")

	c.Srv.Config.LocationName = nam
	c.Srv.Config.Longitude = float32(a)
	c.Srv.Config.Latitude = float32(b)
	c.Srv.Config.Provider = p
	c.Srv.Config.AppID = app
	c.Srv.Config.SetDefaults()

	c.Srv.Config.WriteToFile("config.json")
}

// LogInfo is used to log information messages for this controller.
func (c *ConfigController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Info("ConfigController: ", a[1:len(a)-1])
}

// LogError is used to log error messages for this controller.
func (c *ConfigController) LogError(v ...interface{}) {
	a := fmt.Sprint(v)
	logger.Error("ConfigController: ", a[1:len(a)-1])
}
