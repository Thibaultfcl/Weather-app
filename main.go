package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
}

// constant for the redirect
const port = ":8080"

// index page
func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/index.html")
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	apiKey := os.Getenv("API_KEY")
	apiUrl := "https://api.openweathermap.org/data/2.5/weather?units=metric&q=" + city + "&appid=" + apiKey

	resp, err := http.Get(apiUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(weatherData)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	http.HandleFunc("/", Index)
	http.HandleFunc("/weather", getWeather)

	//load the CSS and the images
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	//start the local host
	fmt.Println("\n(http://localhost:8080/) - Server started on port", port)
	http.ListenAndServe(port, nil)
}
