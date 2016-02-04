package places

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type PlacesClient struct {
	Client      *http.Client
	apiKey      string
	apiEndPoint string
}

type Result struct {
	NextToken string     `json:"next_page_token"`
	Locations []Location `json:"results"`
	Status    string     `json:"status"`
}

type Location struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Vicinity string `json:"vicinity"`

	Hours struct {
		OpenNow bool `json:"open_now"`
	} `json:"opening_hours"`

	Geometry struct {
		Location struct {
			Latitude  float64 `json:"lat"`
			Longitude float64 `json:"lng"`
		}
	} `json:"geometry"`
}

const DefaultApiEndpoint = "https://maps.googleapis.com/maps/api/place"

func NewPlacesClient(apiKey string) *PlacesClient {
	return &PlacesClient{
		apiKey:      apiKey,
		Client:      http.DefaultClient,
		apiEndPoint: DefaultApiEndpoint,
	}
}

func (c *PlacesClient) Nearby(lat float64, lng float64, types ...string) ([]byte, error) {

	latStr := fmt.Sprintf("%.6f", lat)
	lngStr := fmt.Sprintf("%.6f", lng)

	params := make(map[string]string)
	params["location"] = strings.Join([]string{latStr, lngStr}, ",")
	params["rankby"] = "distance"

	if len(types) > 0 {
		params["types"] = strings.Join(types, "|")
	}

	return c.nearBy(params)
}

func (c *PlacesClient) PopularNearby(lat float64, lng float64, radius int, types ...string) ([]byte, error) {
	latStr := fmt.Sprintf("%.6f", lat)
	lngStr := fmt.Sprintf("%.6f", lng)

	params := make(map[string]string)
	params["location"] = strings.Join([]string{latStr, lngStr}, ",")

	params["radius"] = strconv.Itoa(radius)

	params["rankby"] = "prominence"

	if len(types) > 0 {
		params["types"] = strings.Join(types, "|")
	}

	return c.nearBy(params)

}

func (c *PlacesClient) NearbyWithToken(token string) ([]byte, error) {

	params := make(map[string]string)
	params["pagetoken"] = token

	return c.nearBy(params)
}

func (c *PlacesClient) nearBy(params map[string]string) ([]byte, error) {
	return c.dispatchRequest("nearbysearch", params)
}

func (c *PlacesClient) dispatchRequest(reqEndPoint string, params map[string]string) ([]byte, error) {

	reqUrl := strings.Join([]string{c.apiEndPoint, reqEndPoint, "json"}, "/")

	values := url.Values{}
	values.Set("sensor", "true")
	values.Set("key", c.apiKey)

	for k, v := range params {
		values.Set(k, v)
	}

	reqUrl = strings.Join([]string{reqUrl, values.Encode()}, "?")

	req, err := c.Client.Get(reqUrl)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)

	if req.StatusCode >= 200 && req.StatusCode <= 400 && err == nil {
		return body, nil
	} else {
		return nil, fmt.Errorf("Code:%d error:%v body:%s", req.StatusCode, err, body)
	}

}
