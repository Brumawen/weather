package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/IvanMenshykov/MoonPhase"
)

// Moon holds the details about the phase of the moon
type Moon struct {
	Date         time.Time
	Age          float32
	Phase        float32
	PhaseName    string
	Illumination float32
}

// ForDate loads the struct with information about the moon phase for the specified date
func (m *Moon) ForDate(t time.Time) {
	p := MoonPhase.New(t)
	m.Date = t
	m.Age = float32(p.Age())
	m.Phase = float32(p.Phase())
	m.PhaseName = p.PhaseName()
	m.Illumination = float32(p.Illumination())
}

// WriteTo serializes the entity and writes it to the http response
func (m *Moon) WriteTo(w http.ResponseWriter) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	w.Header().Set("content-type", "application/json")
	w.Write(b)
	return nil
}
