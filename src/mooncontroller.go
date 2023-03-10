package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// MoonController handles the Web Methods for retrieving weather and forecast information.
type MoonController struct {
	Srv *Server
}

// AddController adds the controller routes to the router
func (c *MoonController) AddController(router *mux.Router, s *Server) {
	c.Srv = s
	router.Methods("GET").Path("/moon/get").Name("GetMoonCurrent").
		Handler(Logger(c, http.HandlerFunc(c.handleGetCurrent)))
}

// LogInfo is used to log information messages for this controller.
func (c *MoonController) LogInfo(v ...interface{}) {
	a := fmt.Sprint(v...)
	logger.Info("MoonController: ", a)
}

// Get the current weather information
func (c *MoonController) handleGetCurrent(w http.ResponseWriter, r *http.Request) {
	m := Moon{}
	m.ForDate(time.Now())
	if err := m.WriteTo(w); err != nil {
		http.Error(w, "Error serializing moon information. "+err.Error(), 500)
	}
}
