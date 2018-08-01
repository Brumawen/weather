package main

// WeatherProvider provides an interface for a weather information provider
type WeatherProvider interface {
	GetProviderName() string
	GetWeather() (Weather, error)
	GetForecast() (Forecast, error)
	SetConfig(c *Config)
}
