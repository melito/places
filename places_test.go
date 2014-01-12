package places

import (
	"os"
	"testing"
	"time"
)

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

	status := results.(map[string]interface{})["status"]
	if status != "OK" {
		t.Fatal("Request failed.  Got status:", status)
	}

	places := results.(map[string]interface{})["results"]
	if len(places.([]interface{})) == 0 {
		t.Fatal("Didn't get any places:", places)
	}

}

func TestNearbySearchPagination(t *testing.T) {

	client := NewPlacesClient(os.Getenv("PLACES_API_KEY"))
	latitude, longitude := 30.0279, -98.1179

	results, err := client.Nearby(latitude, longitude, "food", "cafe")
	if err != nil {
		t.Fatal("Error returned:", err)
	}

	token := results.(map[string]interface{})["next_page_token"]
	if token == "" {
		t.Fatal("We didn't get a next page token")
	}

	// There is an issue with the page token approach
	// When Google issues you a set of results they provide a token to access the next page.
	// The problem is that that token isn't valid until some time has elapsed.
	time.Sleep(1200 * time.Millisecond)

	nextResults, err := client.NearbyWithToken(token.(string))
	if err != nil {
		t.Fatal("Error returned:", err)
	}

	if nextResults == nil {
		t.Fatal("nextResults was nil")
	}

	status := nextResults.(map[string]interface{})["status"]
	if status != "OK" {
		t.Fatal("Request failed.  Got status:", status)
	}

	places := nextResults.(map[string]interface{})["results"]
	if len(places.([]interface{})) == 0 {
		t.Fatal("Didn't get any places:", places)
	}

}
