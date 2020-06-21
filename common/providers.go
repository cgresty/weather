package common

type WeatherProvider interface {
	OneDayForecast(dayDelta int) OneDayForecast
}

var providers = make(map[string]WeatherProvider)
var defaultProvider WeatherProvider

func AddProvider(name string, provider WeatherProvider) {
	providers[name] = provider
	if len(providers) == 1 {
		defaultProvider = provider
	}
}

func Provider(name string) WeatherProvider {
	return providers[name]
}

func DefaultProvider() WeatherProvider {
	return defaultProvider
}