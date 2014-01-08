package places

import (
	"encoding/json"
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

const DefaultApiEndpoint = "https://maps.googleapis.com/maps/api/place"

func NewPlacesClient(apiKey string) *PlacesClient {
	return &PlacesClient{
		apiKey:      apiKey,
		Client:      http.DefaultClient,
		apiEndPoint: DefaultApiEndpoint,
	}
}

func (c *PlacesClient) Nearby(lat float64, lng float64, radius int, types ...string) (interface{}, error) {

	latStr := fmt.Sprintf("%.6f", lat)
	lngStr := fmt.Sprintf("%.6f", lng)

	params := make(map[string]string)
	params["location"] = strings.Join([]string{latStr, lngStr}, ",")

	params["radius"] = strconv.Itoa(radius)

	if len(types) > 0 {
		params["types"] = strings.Join(types, "|")
	}

	return c.nearBy(params)
}

func (c *PlacesClient) nearBy(params map[string]string) (interface{}, error) {
	return c.dispatchRequest("nearbysearch", params)
}

func (c *PlacesClient) dispatchRequest(reqEndPoint string, params map[string]string) (interface{}, error) {

	reqUrl := strings.Join([]string{c.apiEndPoint, reqEndPoint, "json"}, "/")

	values := url.Values{}
	values.Set("key", c.apiKey)
	values.Set("sensor", "true")

	for k, v := range params {
		values.Set(k, v)
	}

	reqUrl = strings.Join([]string{reqUrl, values.Encode()}, "?")

	req, err := c.Client.Get(reqUrl)
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)

	if req.StatusCode >= 200 && req.StatusCode <= 400 && err == nil {
		var jsonResult interface{}
		json.Unmarshal(body, &jsonResult)
		return jsonResult, nil
	} else {
		return nil, fmt.Errorf("Code:%d error:%v body:%s", req.StatusCode, err, body)
	}

}