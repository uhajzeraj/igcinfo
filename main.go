package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	igc "github.com/marni/goigc"
)

// Some .igc files URLs
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Jarez%20to%20Senegal.igc
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Boavista%20Medellin.igc
// http://skypolaris.org/wp-content/uploads/IGS%20Files/Medellin%20Guatemala.igc

const urlRoot = "/igcinfo/"

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

// Get the index of the track in the igcFiles slice, if it is there
func getTrackIndex(uniqueID string) int {
	for index, trackInArray := range igcFiles {
		if trackInArray.UniqueID == uniqueID {
			return index
		}
	}
	return -1
}

// ISO8601 duration parsing function
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

// Calculate the total distance of the track
func calculateTotalDistance(track igc.Track) string {

	totalDistance := 0.0

	// For each point of the track, calculate the distance between 2 points in the Point array
	for i := 0; i < len(track.Points)-1; i++ {
		totalDistance += track.Points[i].Distance(track.Points[i+1])
	}

	// Parse it to a string value
	return strconv.FormatFloat(totalDistance, 'f', 2, 64)
}

// Check if any of the regex patterns supplied in the map parameter match the string parameter
func regexMatches(url string, urlMap map[string]func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	for mapURL := range urlMap {
		res, err := regexp.MatchString(mapURL, url)
		if err != nil {
			return nil
		}

		if res { // If the pattern matching returns true, return the function
			return urlMap[mapURL]
		}
	}
	return nil
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

func apiIgcIDHandler(w http.ResponseWriter, r *http.Request) {
	urlID := path.Base(r.URL.Path) // returns the part after the last '/' in the url

	trackSliceID := getTrackIndex(urlID)
	if trackSliceID != -1 { // Check whether the ID is different from -1
		w.Header().Set("Content-Type", "application/json") // Set response content-type to JSON

		response := "{"
		response += "\"H_date\": " + "\"" + igcFiles[trackSliceID].Date.String() + "\","
		response += "\"pilot\": " + "\"" + igcFiles[trackSliceID].Pilot + "\","
		response += "\"glider\": " + "\"" + igcFiles[trackSliceID].GliderType + "\","
		response += "\"glider_id\": " + "\"" + igcFiles[trackSliceID].GliderID + "\","
		response += "\"track_length\": " + "\"" + calculateTotalDistance(igcFiles[trackSliceID]) + "\"" // TO-DO, calculate the track length?
		response += "}"

		fmt.Fprintf(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

func apiIgcIDFieldHandler(w http.ResponseWriter, r *http.Request) {

	pathArray := strings.Split(r.URL.Path, "/") // split the URL Path into chunks, whenever there's a "/"
	field := pathArray[len(pathArray)-1]        // The part after the last "/", is the field
	uniqueID := pathArray[len(pathArray)-2]     // The part after the second to last "/", is the unique ID

	trackSliceID := getTrackIndex(uniqueID)

	if trackSliceID != -1 { // Check whether the ID is different from -1

		something := map[string]string{ // Map the field to one of the Track struct attributes in the igcFiles slice
			"pilot":        igcFiles[trackSliceID].Pilot,
			"glider":       igcFiles[trackSliceID].GliderType,
			"glider_id":    igcFiles[trackSliceID].GliderID,
			"track_length": calculateTotalDistance(igcFiles[trackSliceID]),
			"H_date":       igcFiles[trackSliceID].Date.String(),
		}

		response := something[field] // This will work because the RegEx checks whether the name is written correctly
		fmt.Fprintf(w, response)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

func urlRouter(w http.ResponseWriter, r *http.Request) {

	urlMap := map[string]func(http.ResponseWriter, *http.Request){ // A map of accepted URL RegEx patterns
		"^/igcinfo/api/$":                      apiHandler,
		"^/igcinfo/api/igc$":                   apiIgcHandler,
		"^/igcinfo/api/igc/[a-zA-Z0-9]{3,10}$": apiIgcIDHandler,
		"^/igcinfo/api/igc/[a-zA-Z0-9]{3,10}/(pilot|glider|glider_id|track_length|H_date)$": apiIgcIDFieldHandler,
	}

	result := regexMatches(r.URL.Path, urlMap) // Perform the RegEx check to see if any pattern matches

	if result != nil { // If a function is returned, call that handler function
		result(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound) // If it isn't, send a 404 Not Found status
	}
}

func main() {
	http.HandleFunc("/", urlRouter) // Handle all the request via the urlRouter function
	log.Fatal(http.ListenAndServe(":8080", nil))
}
