package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	igc "github.com/marni/goigc"
)

// Some .igc files URLs
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc

// Slice where the igcFiles are in-memory stored
var igcFiles []igc.Track

// Unix timestamp when the service started
var timeStarted = int(time.Now().Unix())

// Check if uniqueID is in the igcFiles slice
func stringInSlice(uniqueID string) bool {
	for _, trackInArray := range igcFiles {
		if trackInArray.UniqueID == uniqueID {
			return true
		}
	}
	return false
}

func stringInMap(url string, urlMap map[string]func(http.ResponseWriter, *http.Request)) bool {
	for mapURL := range urlMap {
		if mapURL == url {
			return true
		}
	}

	return false
}

func apiIgcHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" { // If method is POST, user has entered the URL
		var data map[string]string // POST body is of content-type: JSON; the result can be stored in a map
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			panic(err)
		}

		track, err := igc.ParseLocation(data["url"]) // call the igc library
		if err != nil {
			panic(err)
		}

		// Check if track slice contains the uniqueID

		if len(igcFiles) != 0 { // In case the slice is empty, just add the track
			if !stringInSlice(track.UniqueID) { // If the uniqueID is not in the slice, add it
				igcFiles = append(igcFiles, track) // Append the result to igcFiles slice
			}

		} else {
			igcFiles = append(igcFiles, track) // Append the result to igcFiles slice
		}

		response := "{"
		response += "\"id\": " + "\"" + track.UniqueID + "\""
		response += "}"

		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON
		fmt.Fprintf(w, response)

	} else if r.Method == "GET" { // If the method is GET
		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		response := "["
		for i := range igcFiles { // Get all the IDs of .igc files stored in the igcFiles array
			if i != len(igcFiles)-1 { // If it's the last item in the array, don't add the ","
				response += "\"" + igcFiles[i].UniqueID + "\","
			} else {
				response += "\"" + igcFiles[i].UniqueID + "\""
			}
		}
		response += "]"

		fmt.Fprintf(w, response)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

	timeNow := int(time.Now().Unix()) // Unix timestamp when the handler was called

	iso8601duration := parseTimeDifference(timeNow - timeStarted) // Calculate the time elapsed by subtracting the times

	response := "{"
	response += "\"uptime\": \"" + iso8601duration + "\","
	response += "\"info\": \"Service for IGC tracks.\","
	response += "\"version\": \"v1\""
	response += "}"
	fmt.Fprintln(w, response)
}

func parseTimeDifference(timeDifference int) string {

	result := "P" // Different time intervals are attached to this, if they are != 0

	// Formulas for calculating different time intervals in seconds
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

	if hours != 0 || minutes != 0 || seconds != 0 { // Check in case time intervals are 0
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

func urlRouter(w http.ResponseWriter, r *http.Request) {

	// **TO DO** Change the static URLs to RegEx patterns
	urlMap := map[string]func(http.ResponseWriter, *http.Request){ // A map of accepted URLs
		"/igcinfo/api/":    apiHandler,
		"/igcinfo/api/igc": apiIgcHandler,
	}

	if stringInMap(r.URL.Path, urlMap) { // Check if the request is in the map
		urlMap[r.URL.Path](w, r) // If it is, redirect to that handler
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

// w.WriteHeader(status)

func main() {
	http.HandleFunc("/", urlRouter)
	// http.HandleFunc("/igcinfo/api/", apiHandler)
	// http.HandleFunc("/igcinfo/api/igc", apiIgcHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
