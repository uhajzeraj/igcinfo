package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	igc "github.com/marni/goigc"
)

// Unix timestamp when the service started
var timeStarted = int(time.Now().Unix())

type url struct {
	URL string `json:"url"`
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("igc.html")

	if r.Method == "POST" {

		apiURL := &url{}
		apiURL.URL = r.FormValue("url")
		// var jsonR map[string]string
		var _ = json.NewDecoder(r.Body).Decode(apiURL)

		track, _ := igc.ParseLocation(apiURL.URL)

		response := "{"
		response += "\"id\": " + "\"" + track.UniqueID + "\","
		response += "\"url\": " + "\"" + apiURL.URL + "\""
		response += "}"

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, response)

	} else if r.Method == "GET" {
		t.Execute(w, nil)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	// Set response content-type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Unix timestamp when the handler was called
	timeNow := int(time.Now().Unix())

	iso8601duration := parseTimeDifference(timeNow - timeStarted)

	// Calculate the time elapsed by subtracting the times
	response := "{"
	response += "\"uptime\": \"" + iso8601duration + "\","
	response += "\"info\": \"Service for IGC tracks.\","
	response += "\"version\": \"v1\""
	response += "}"
	fmt.Fprintln(w, response)
}

func parseTimeDifference(timeDifference int) string {

	result := "P" // Different time intervals are attached to this, if they are != 0

	// Formulas for calculating different time intervals
	timeLeft := timeDifference
	years := timeDifference / 31557600
	timeLeft -= years * 31557600
	months := timeLeft / 2592000
	timeLeft -= months * 2592000
	weeks := timeLeft / 604800
	timeLeft -= weeks * 604800
	days := timeLeft / 86400
	timeLeft -= days * 86400
	hours := timeLeft / 3600
	timeLeft -= hours * 3600
	minutes := timeLeft / 60
	timeLeft -= minutes * 60
	seconds := timeLeft

	// Add time invervals to the result only if they are different form 0
	if years != 0 {
		result += fmt.Sprintf("Y%d", years)
	}
	if months != 0 {
		result += fmt.Sprintf("M%d", months)
	}
	if weeks != 0 {
		result += fmt.Sprintf("W%d", weeks)
	}
	if days != 0 {
		result += fmt.Sprintf("D%d", days)
	}

	// Check in case time intervals are 0
	if hours != 0 || minutes != 0 || seconds != 0 {
		result += "T"
		if hours != 0 {
			result += fmt.Sprintf("H%d", hours)
		}
		if minutes != 0 {
			result += fmt.Sprintf("M%d", minutes)
		}
		if seconds != 0 {
			result += fmt.Sprintf("S%d", seconds)
		}
	}

	return result
}

func main() {
	http.HandleFunc("/igcinfo/api/", apiHandler)
	http.HandleFunc("/igcinfo/api/igc", postHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
