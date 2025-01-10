package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	baseURL    = "https://api.open-meteo.com/v1/forecast"
	geocodeURL = "https://geocoding-api.open-meteo.com/v1/search"
)

type WeatherResponse struct {
	Daily struct {
		Time           []string  `json:"time"`
		TemperatureMax []float64 `json:"temperature_2m_max"`
		TemperatureMin []float64 `json:"temperature_2m_min"`
	} `json:"daily"`
}

type GeocodeResponse struct {
	Results []struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"results"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run weather.go <CITY_NAME>")
		return
	}

	city := os.Args[1]

	// Get coordinates from city name
	latitude, longitude, err := getCoordinates(city)
	if err != nil {
		fmt.Printf("Failed to get coordinates: %v\n", err)
		return
	}

	// Fetch 5-day weather forecast
	forecastURL := fmt.Sprintf("%s?latitude=%f&longitude=%f&daily=temperature_2m_max,temperature_2m_min&timezone=auto", baseURL, latitude, longitude)
	weatherData, err := fetchForecast(forecastURL)
	if err != nil {
		fmt.Printf("Failed to get weather data: %v\n", err)
		return
	}

	// Output the forecast
	fmt.Printf("Weather forecast for %s:\n", city)
	display5DayForecast(weatherData)
}

func getCoordinates(city string) (float64, float64, error) {
	url := fmt.Sprintf("%s?name=%s", geocodeURL, city)
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return 0, 0, fmt.Errorf("Error: %s\nResponse Body: %s", resp.Status, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var geocodeData GeocodeResponse
	if err := json.Unmarshal(body, &geocodeData); err != nil {
		return 0, 0, err
	}

	if len(geocodeData.Results) == 0 {
		return 0, 0, fmt.Errorf("No results found for city: %s", city)
	}

	return geocodeData.Results[0].Latitude, geocodeData.Results[0].Longitude, nil
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
	for i, date := range data.Daily.Time {
		// Parse the date to get the day of the week
		parsedDate, err := time.Parse("2006-01-02", date)
		if err != nil {
			fmt.Printf("Failed to parse date: %v\n", err)
			continue
		}
		dayOfWeek := parsedDate.Weekday()
		fmt.Printf("%s (%s) - High: %.1f°C, Low: %.1f°C\n", dayOfWeek, date, data.Daily.TemperatureMax[i], data.Daily.TemperatureMin[i])
	}
}
