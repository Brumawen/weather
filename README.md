# weather

Weather is a microservice that can be used to get the current weather and forecast for a particular location.


## Installation

Download the latest release file for your operating system from https://github.com/Brumawen/weather/releases 

Extract the files and subfolders to a folder and run the following from the command line

        weather -service install
        weather -service run

This will install and run the weather microservice as a background service on your machine.


## Configuration

Once the microservice is running, navigate to http://localhost:20511/config.html in a web browser.

The service will automatically detect the location based on your machine's public IP address.  You can change the Location information to be more accurate for your location.

You will need to obtain an Application ID for the chosen weather provider.  The following providers are available:

* Open Weather (https://openweathermap.org/appid)
* AccuWeather (https://developer.accuweather.com/)

Paste your APPID value into the Application ID field and click Save.

## Weather Display

To display the current weather and forecast details, navifate to http://localhost:20511/weather.html


# Weather API

To get the current weather information for the configured location.

        http://localhost:20511/weather/current

* Created: The date and time the information was collated by the weather microservice.
* Humidity: The current humidity (%).
* ID: The location identifier.
* IsDay: Returns true if the current time is day time.
* Name: The location name
* Pressure: The current pressure (mb or inHg).
* Provider:  Name of the weather provider.
* ReadingTime: The date and time the reading was taken by the provider.
* Sunrise: The time of sunrise.
* Sunset: The time of sunset.
* Temp: The current temperature (celcius or farenheit).
* WeatherIcon: The icon to use for the weather.  See weather icons below.
* WeatherDesc: Weather description.
* WindSpeed: The wind speed in (km/h or m/h).
* WindDirection: The cardinal direction the wind is coming from.

To get the 5 day weather forecast

        http://localhost:20511/weather/forecast

* Day: The date of the forecast
* Name: The name of the day
* TempMin: The minumum expected temperature (celcius or farenheit)
* TempMax: The maximum expected temperature (celcius or farenheit)
* WeatherIcon: The icon to use for the weather.  See weather icons below.
* WeatherDesc: Weather description.

Weather Icons

1. Sunny
2. Scattered Clouds
3. Partly Cloudy
4. Cloudy
5. Scattered Rain
6. Rain
7. Thunderstorms
8. Snow
9. Mist

## Moon Phase API

To get the current phase of the moon

        http://localhost:20511/moon/get

* Date: Date of the moon phase.
* Age: Age of the moon (0 to 28 days).
* Phase: Phase as a value from 0 (new) to 1 (full).
* PhaseName: Name of the phase.
* Illumination: Amount of illumination from 0 (new) to 1 (full). 