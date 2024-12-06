package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := os.ReadFile(filename)

	var c apiConfigData

	if err != nil {
		return c, err
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/{location}",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(data)
		})

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello from go!\n"))
}

func query(city string) (weatherData, error) {
	var data weatherData

	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return data, err
	}

	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city)
	if err != nil {
		return data, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return data, nil
	}

	return data, nil
}
