package main

import (
	"fmt"
	"net/http"
)

// Serve the test HTML page
func handler(w http.ResponseWriter, r *http.Request) {
	// Define a simple HTML page with embedded CSS and JavaScript
	html := `<!DOCTYPE html>
	<html lang="en">
	<head>
	  <meta charset="UTF-8">
	  <meta name="viewport" content="width=device-width, initial-scale=1.0">
	  <title>Weather Forecast</title>
	  <!-- Material UI CSS -->
	  <link href="https://fonts.googleapis.com/css2?family=Roboto:wght@400;500&display=swap" rel="stylesheet">
	  <link href="https://cdn.jsdelivr.net/npm/@mui/material@5.10.0/dist/material-ui.min.css" rel="stylesheet">
	  <style>
	    body {
	      font-family: 'Roboto', sans-serif;
	      background-color: #f4f4f9;
	      margin: 0;
	      padding: 0;
	    }
	    .container {
	      max-width: 800px;
	      margin: 50px auto;
	      padding: 20px;
	      background-color: white;
	      border-radius: 8px;
	      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	    }
	    h1 {
	      text-align: center;
	      color: #333;
	    }
	    .location {
	      text-align: center;
	      font-size: 1.2rem;
	      color: #555;
	      margin-top: 10px;
	    }
	    .forecast {
	      display: grid;
	      grid-template-columns: repeat(5, 1fr);
	      gap: 20px;
	    }
	    .forecast-item {
	      background-color: #e3f2fd;
	      border-radius: 8px;
	      padding: 15px;
	      text-align: center;
	    }
	    #loading {
	      text-align: center;
	      font-size: 1.2rem;
	      color: #00796b;
	    }
	  </style>
	</head>
	<body>
	  <div class="container">
	    <h1>Weather Forecast for the Next 5 Days</h1>
	    <div id="location" class="location"></div> <!-- Display geolocated location -->
	    <div id="loading">Loading...</div>
	    <div id="forecast" class="forecast"></div>
	  </div>

	  <!-- React and Material UI JS -->
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/react/17.0.2/umd/react.production.min.js"></script>
	  <script src="https://cdnjs.cloudflare.com/ajax/libs/react-dom/17.0.2/umd/react-dom.production.min.js"></script>
	  <script src="https://cdn.jsdelivr.net/npm/@mui/material@5.10.0/dist/material-ui.min.js"></script>

	  <script>
  const apiKey = '270e52f7980a4cbe870231539240112'; // Replace with your WeatherAPI key

  // Function to fetch weather based on geolocation
function getWeather(latitude, longitude) {
  const url = 'https://api.weatherapi.com/v1/forecast.json?key=' + apiKey + '&q=' + latitude + ',' + longitude + '&days=5&aqi=no&alerts=no';
  fetch(url)
    .then(response => response.json())
    .then(data => {
      const forecast = data.forecast.forecastday;
      const forecastContainer = document.getElementById('forecast');
      const locationElement = document.getElementById('location');

      // Set the location name (optional)
      locationElement.innerText = 'Weather from: ' + data.location.name + ', ' + data.location.region;

      forecastContainer.innerHTML = ''; // Clear loading message

      forecast.forEach(function(day) {
        const dayDiv = document.createElement('div');
        dayDiv.className = 'forecast-item';
        dayDiv.innerHTML = '<h3>' + new Date(day.date).toLocaleDateString() + '</h3>' +
                          '<p>' + day.day.condition.text + '</p>' +
                          '<p>High: ' + day.day.maxtemp_c + '&#176;C</p>' +  // Use `&#176;` for the degree symbol
                          '<p>Low: ' + day.day.mintemp_c + '&#176;C</p>';
        forecastContainer.appendChild(dayDiv);
      });
    })
    .catch(function(error) {
      console.error('Error fetching weather data:', error);
      document.getElementById('loading').innerText = 'Failed to load weather data';
    });
}

// Function to get the user's geolocation
function getGeolocation() {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      function(position) {
        const latitude = position.coords.latitude;
        const longitude = position.coords.longitude;
        getWeather(latitude, longitude);
      },
      function(error) {
        console.error('Error getting geolocation:', error);
        document.getElementById('loading').innerText = 'Failed to retrieve your location';
      }
    );
  } else {
    console.error('Geolocation is not supported by this browser.');
    document.getElementById('loading').innerText = 'Geolocation not supported';
  }
}

// Start the process of getting geolocation and weather
window.onload = function() {
  getGeolocation();
};
</script>

	</body>
	</html>`

	// Write the HTML content to the response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, html)
}

func main() {
	// Handle requests to the root URL by serving the handler function
	http.HandleFunc("/", handler)

	// Start the server on port 8080
	fmt.Println("Starting server on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
