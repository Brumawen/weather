# weather

Weather is a microservice that can be used to get the current weather and forecast for a particular location.


## Installation

Download the release file for your operating system from https://github.com/Brumawen/weather/releases 

Extract the file to a folder and run the following from the command line

        weather -server install
        weather -server run

This will install and run the weather microservice as a background service on your machine.


## Configuration

Once the microservice is running, navigate to http://localhost:20511/config.html in a web browser.

The service will automatically detect the location based on your machine's public IP address.  You can change the Location information to be more accurate for your location.

Currently the only supported provider is Open Weather.  You will need to obtain an Application ID from here https://openweathermap.org/appid

Paste your APPID value into the Application ID field and click Save.


## Weather

To get the current weather

        http://localhost:20511/weather/current

To get the 5 day weather forecast

        http://localhost:20511/weather/forecast


## Moon Phase

To get the current phase of the moon

        http://localhost:20511/moon/get

