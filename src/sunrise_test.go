package main

import "testing"
import "time"
import "fmt"

func TestCanGetSunriseSunset(t *testing.T) {
	c := &Config{}
	err := c.ReadFromFile("config.json")
	if err != nil {
		t.Error(err)
	}

	sr, ss, err := GetSunriseSunset(c, time.Now())
	if err != nil {
		t.Error(err)
	}
	fmt.Println("Sunrise =", sr.String())
	fmt.Println("Sunset =", ss.String())
}
