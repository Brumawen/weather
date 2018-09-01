package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// IPLocation holds the information return from a call to get the location information
// of the public IP address of the current computer
type IPLocation struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Latitude    float32 `json:"lat"`
	Longitude   float32 `json:"lon"`
	Region      string  `json:"regionName"`
	RegionCode  string  `json:"region"`
	TimeZone    string  `json:"timezone"`
	IPAddress   string  `json:"query"`
}

// GetPublicIPAddress returns the public IP address of the current computer.
func GetPublicIPAddress() (string, error) {
	resp, err := http.Get("http://checkip.amazonaws.com/")
	if resp != nil {
		defer resp.Body.Close()
		resp.Close = true
	}
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), err
}

// GetIPLocationInfo gets the location information for the current computer based on
// the public IP address
func GetIPLocationInfo() (IPLocation, error) {
	info := IPLocation{}
	resp, err := http.Get("http://ip-api.com/json")
	if resp != nil {
		defer resp.Body.Close()
		resp.Close = true
	}
	if err != nil {
		return info, err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return info, err
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		return info, err
	}
	return info, nil
}
