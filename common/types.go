package common

import "time"

type Hour int

type Forecast struct {
	Time time.Time
	Location string
	Temperature int
	WindDirection int
	WindSpeedMph int
	WindGustMph int
	WeatherType int
}

type OneDayForecast struct {
	Time time.Time
	Headline string
	Description string
	Forecasts map[Hour]Forecast
}
