package places

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func structureResult(data []byte) (result *Result) {
	var jsonObj Result
	json.Unmarshal(data, &jsonObj)
	return &jsonObj
}

func TestClientCreation(t *testing.T) {

	client := NewPlacesClient(os.Getenv("PLACES_API_KEY"))

	if client.apiKey == "" {
		t.Fatal("Didn't set an api key")
	}

	if client.apiEndPoint == "" {
		t.Fatal("Didn't set an api endpoint")
	}

}

func TestNearbySearch(t *testing.T) {

	client := NewPlacesClient(os.Getenv("PLACES_API_KEY"))
	latitude, longitude := 30.0279, -98.1179

	results, err := client.Nearby(latitude, longitude, "food", "cafe")
	if err != nil {
		t.Fatal("Error returned:", err)
	}

	if results == nil {
		t.Fatal("Results were nil")
	}

	result := structureResult(results)
	if result.Status != "OK" {
		t.Fatal("Request failed.  Got status:", result.Status)
	}

	if len(result.Locations) == 0 {
		t.Fatal("Didn't get any places:", result.Locations)
	}

}

func TestNearbySearchPagination(t *testing.T) {

	client := NewPlacesClient(os.Getenv("PLACES_API_KEY"))
	latitude, longitude := 30.0279, -98.1179

	results, err := client.Nearby(latitude, longitude, "food", "cafe")
	if err != nil {
		t.Fatal("Error returned:", err)
	}

	result := structureResult(results)

	if result.NextToken == "" {
		t.Fatal("We didn't get a next page token")
	}

	// There is an issue with the page token approach
	// When Google issues you a set of results they provide a token to access the next page.
	// The problem is that that token isn't valid until some time has elapsed.
	time.Sleep(1200 * time.Millisecond)

	nextResults, err := client.NearbyWithToken(result.NextToken)
	if err != nil {
		t.Fatal("Error returned:", err)
	}

	if nextResults == nil {
		t.Fatal("nextResults was nil")
	}

	nextResult := structureResult(nextResults)

	t.Log(nextResult)

	if nextResult.Status != "OK" {
		t.Fatal("Request failed.  Got status:", nextResult.Status)
	}

	if len(nextResult.Locations) == 0 {
		t.Fatal("Didn't get any places:", nextResult.Locations)
	}

}
