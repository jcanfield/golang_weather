package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const apiKey = "270e52f7980a4cbe870231539240112" // Replace with your actual WeatherAPI.com API key
const baseURL = "https://api.weatherapi.com/v1"

type LocationResponse struct {
	Location struct {
		Name      string `json:"name"`
		Region    string `json:"region"`
		Country   string `json:"country"`
	} `json:"location"`
}

type WeatherResponse struct {
	Forecast struct {
		Forecastday []struct {
			Date string `json:"date"`
			Day  struct {
				MaxtempF float64 `json:"maxtemp_f"`
				MintempF float64 `json:"mintemp_f"`
			} `json:"day"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run weather.go <ZIP_CODE>")
		return
	}

	zipCode := os.Args[1]

	// Step 1: Fetch location information (city and state)
	locationURL := fmt.Sprintf("%s/current.json?key=%s&q=%s", baseURL, apiKey, zipCode)
	locationData, err := fetchLocation(locationURL)
	if err != nil {
		fmt.Printf("Failed to get location data: %v\n", err)
		return
	}

	// Step 2: Fetch 5-day weather forecast
	forecastURL := fmt.Sprintf("%s/forecast.json?key=%s&q=%s&days=5", baseURL, apiKey, zipCode)
	weatherData, err := fetchForecast(forecastURL)
	if err != nil {
		fmt.Printf("Failed to get weather data: %v\n", err)
		return
	}

	// Output the forecast with the city and state
	fmt.Printf("Weather forecast for %s, %s:\n", locationData.Location.Name, locationData.Location.Region)
	display5DayForecast(weatherData)
}

func fetchLocation(url string) (*LocationResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error: %s\nResponse Body: %s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var locationData LocationResponse
	if err := json.Unmarshal(body, &locationData); err != nil {
		return nil, err
	}
	return &locationData, nil
}

func fetchForecast(url string) (*WeatherResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error: %s\nResponse Body: %s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherData WeatherResponse
	if err := json.Unmarshal(body, &weatherData); err != nil {
		return nil, err
	}
	return &weatherData, nil
}

func display5DayForecast(data *WeatherResponse) {
	fmt.Println("\n5-Day Forecast:")
	for _, day := range data.Forecast.Forecastday {
		// Parse the date to get the day of the week
		date, err := time.Parse("2006-01-02", day.Date)
		if err != nil {
			fmt.Printf("Failed to parse date: %v\n", err)
			continue
		}
		dayOfWeek := date.Weekday()
		fmt.Printf("%s (%s) - High: %.1f°F, Low: %.1f°F\n", dayOfWeek, day.Date, day.Day.MaxtempF, day.Day.MintempF)
	}
}
