package simworks

// Environment represents the simulation environment
type Environment struct {
	Weather WeatherModel
	// Add other environment-specific attributes here
}

// WeatherModel represents the weather conditions in the simulation
type WeatherModel struct {
	Temperature   float64 // Temperature in Celsius
	Humidity      float64 // Humidity percentage
	WindSpeed     float64 // Wind speed in m/s
	Precipitation float64 // Precipitation in mm/hour
}

// NewEnvironment creates a new environment
func NewEnvironment() *Environment {
	return &Environment{
		Weather: WeatherModel{
			Temperature:   20,
			Humidity:      60,
			WindSpeed:     5,
			Precipitation: 0,
		},
	}
}

// UpdateWeather updates the weather conditions in the environment
func (e *Environment) UpdateWeather(temp, humidity, windSpeed, precip float64) {
	e.Weather = WeatherModel{
		Temperature:   temp,
		Humidity:      humidity,
		WindSpeed:     windSpeed,
		Precipitation: precip,
	}
}
