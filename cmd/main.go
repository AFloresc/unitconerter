package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type ConversionRequest struct {
	Value    float64 `json:"value"`
	FromUnit string  `json:"fromUnit"`
	ToUnit   string  `json:"toUnit"`
}

func convertUnits(value float64, fromUnit, toUnit string) float64 {
	conversionRates := map[string]map[string]float64{
		"millimeter": {"centimeter": 0.1, "meter": 0.001, "kilometer": 0.000001},
		"centimeter": {"millimeter": 10, "meter": 0.01, "kilometer": 0.00001},
		"meter":      {"millimeter": 1000, "centimeter": 100, "kilometer": 0.001},
		"kilometer":  {"millimeter": 1000000, "centimeter": 100000, "meter": 1000},
		"inch":       {"foot": 0.0833, "yard": 0.0278, "mile": 0.0000157},
		"foot":       {"inch": 12, "yard": 0.333, "mile": 0.000189},
		"yard":       {"inch": 36, "foot": 3, "mile": 0.000568},
		"mile":       {"inch": 63360, "foot": 5280, "yard": 1760},

		"milligram": {"gram": 0.001, "kilogram": 0.000001},
		"gram":      {"milligram": 1000, "kilogram": 0.001},
		"kilogram":  {"milligram": 1000000, "gram": 1000},
		"ounce":     {"pound": 0.0625},
		"pound":     {"ounce": 16},

		"Celsius":    {"Fahrenheit": value*9/5 + 32, "Kelvin": value + 273.15},
		"Fahrenheit": {"Celsius": (value - 32) * 5 / 9, "Kelvin": (value-32)*5/9 + 273.15},
		"Kelvin":     {"Celsius": value - 273.15, "Fahrenheit": (value-273.15)*9/5 + 32},
	}

	return conversionRates[fromUnit][toUnit] * value
}

func handleConversion(w http.ResponseWriter, r *http.Request) {
	value, _ := strconv.ParseFloat(r.URL.Query().Get("value"), 64)
	fromUnit := r.URL.Query().Get("fromUnit")
	toUnit := r.URL.Query().Get("toUnit")

	convertedValue := convertUnits(value, fromUnit, toUnit)

	result := map[string]float64{"result": convertedValue}
	json.NewEncoder(w).Encode(result)
}

func main() {
	fs := http.FileServer(http.Dir("webapp"))
	http.Handle("/", fs)
	http.HandleFunc("/convert", handleConversion)
	fmt.Println("Starting server on :8080...")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Unit Converter API!")
}
