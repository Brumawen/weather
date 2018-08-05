package main

import (
	"fmt"
	"time"

	"github.com/kelvins/sunrisesunset"
)

func GetSunriseSunset(c *Config, t time.Time) (time.Time, time.Time, error) {
	y := t.Year()
	m := t.Month()
	d := t.Day()
	n, o := t.Zone()
	fmt.Println(n, o)

	p := sunrisesunset.Parameters{
		Latitude:  float64(c.Latitude),
		Longitude: float64(c.Longitude),
		UtcOffset: float64(o) / 3600,
		Date:      time.Date(y, m, d, 0, 0, 0, 0, time.UTC),
	}
	sr, ss, err := p.GetSunriseSunset()
	if err == nil {
		h := sr.Hour()
		n := sr.Minute()
		s := sr.Second()
		sr = time.Date(y, m, d, h, n, s, 0, time.Local)

		h = ss.Hour()
		n = ss.Minute()
		s = ss.Second()
		ss = time.Date(y, m, d, h, n, s, 0, time.Local)
	}
	return sr, ss, err
}
