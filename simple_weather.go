package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	apiURL = "https://api.weatherapi.com/v1/forecast.json"
	apiKey = "270e52f7980a4cbe870231539240112" // Replace with your WeatherAPI key
)

type WeatherResponse struct {
	Location struct {
		Name string `json:"name"`
	} `json:"location"`
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Day  struct {
				AvgTempF float64 `json:"avgtemp_f"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <zip_code>")
		os.Exit(1)
	}

	zip := os.Args[1]
	url := fmt.Sprintf("%s?key=%s&q=%s&days=4", apiURL, apiKey, zip)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Failed to fetch weather data: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", resp.Status)
		os.Exit(1)
	}

	var weatherData WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		fmt.Printf("Error parsing weather data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Weather for %s:\n", weatherData.Location.Name)
	for _, day := range weatherData.Forecast.Forecastday {
		fmt.Printf("%s: %.2f°F, %s\n", day.Date, day.Day.AvgTempF, day.Day.Condition.Text)
	}
}
