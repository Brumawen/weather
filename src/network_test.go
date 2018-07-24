package main

import (
	"fmt"
	"testing"
)

func TestCanGetPublicIP(t *testing.T) {
	s, err := GetPublicIPAddress()
	if err != nil {
		t.Error(err)
	}
	if s == "" {
		t.Error("The IP address is blank")
	}
	fmt.Println("IP address is", s)

}

func TestCanGetIPLocation(t *testing.T) {
	l, err := GetIPLocationInfo()
	if err != nil {
		t.Error(err)
	}
	if l.IPAddress == "" {
		t.Error("The IP address is blank")
	}
	if l.Latitude == 0 {
		t.Error("The Latitude is 0")
	}
	if l.Longitude == 0 {
		t.Error("The Longitude is 0")
	}
	if l.City == "" {
		t.Error("The City is blank")
	}
	if l.Country == "" {
		t.Error("The Country is blank")
	}
	if l.CountryCode == "" {
		t.Error("The CountryCode is blank")
	}
	fmt.Println(l)
}
