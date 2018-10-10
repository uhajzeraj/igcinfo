package main

import (
	"net/http"
	"testing"
)

func Test_urlRouter(t *testing.T) {

	urlMap := map[string]func(http.ResponseWriter, *http.Request){
		"^/igcinfo/api$":                    apiHandler,
		"^/igcinfo/api/igc$":                apiIgcHandler,
		"^/igcinfo/api/igc/igc[0-9]{1,10}$": apiIgcIDHandler,
		"^/igcinfo/api/igc/[a-zA-Z0-9]{3,10}/(pilot|glider|glider_id|track_length|H_date)$": apiIgcIDFieldHandler,
	}

	// Valid testing cases
	validTestUrls := []string{"/igcinfo/api", "/igcinfo/api/igc", "/igcinfo/api/igc/igc2", "/igcinfo/api/igc/igc2/pilot",
		"/igcinfo/api/igc/igc2/glider", "/igcinfo/api/igc/igc2/glider_id", "/igcinfo/api/igc/igc2/track_length",
		"/igcinfo/api/igc/igc2/H_date"}

	// This function is supposed to return a handlerFunction
	// In case it does not (regexMatches(...) == nil), it means the test has failed
	for _, url := range validTestUrls {
		if regexMatches(url, urlMap) == nil {
			t.Error("Expected apiHandler function, got nothing")
		}
	}

	// Invalid testing cases
	invalidTestUrls := []string{"/", "", "rubbishtext", "/igcinfo/rubbishtext", "/igcinfo/api/igc/somemorerubbish", "/igcinfo/api/igc/igc2/rubish",
		"/igcinfo/api/igc/igc2/track_length+someRubbish", "/igcinfo/api/igc/igc2/H_date/addtionalRubbish"}

	// This function is not supposed to return a handlerFunction
	// In case it does (regexMatches(...) != nil), it means the test has failed
	for _, url := range invalidTestUrls {
		if regexMatches(url, urlMap) != nil {
			t.Error("Unexpected apiHandler function, got something")
		}
	}

}

func Test_parseTimeDifference(t *testing.T) {

	secondsArray := []int{0, 1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576,
		2097152, 4194304, 8388608, 16777216, 33554432, 67108864, 134217728, 268435456, 536870912, 1073741824, 2147483648, 4294967296}

	returnValueArray := []string{"P", "PTS1", "PTS2", "PTS4", "PTS8", "PTS16", "PTS32", "PTM1S4", "PTM2S8", "PTM4S16", "PTM8S32", "PTM17S4", "PTM34S8",
		"PTH1M8S16", "PTH2M16S32", "PTH4M33S4", "PTH9M6S8", "PTH18M12S16", "PD1TH12M24S32", "PD3TM49S4", "PD6TH1M38S8", "PW1D5TH3M16S16",
		"PW3D3TH6M32S32", "PM1W2D4TH13M5S4", "PM3W1TH2M10S8", "PM6W2TH4M20S16", "PY1W3D2TH2M40S32", "PY2M1W2D2TH5M21S4",
		"PY4M3D2TH10M42S8", "PY8M6D4TH21M24S16", "PY17D4TH12M48S32", "PY34W1D2TH1M37S4", "PY68W2D4TH3M14S8", "PY136M1D6TH6M28S16"}

	// Check whether the seconds in secondsArray correspond to the formatted
	for i := 0; i < len(returnValueArray); i++ {
		if parseTimeDifference(secondsArray[i]) != returnValueArray[i] {
			t.Error("Time duration format is not correct")
		}
	}
}
