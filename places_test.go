package places

import (
	"os"
	"testing"
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

	results, err := client.Nearby(latitude, longitude, 5000, "food", "cafe")
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
