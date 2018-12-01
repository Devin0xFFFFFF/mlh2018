package api_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	appID       string
	endpointKey string
	region      string
}

func NewClient(appID string, endpointKey string, region string) *Client {
	return &Client{appID: appID, endpointKey: endpointKey, region: region}
}

func (c *Client) Predict(utterance string) (PredictionResult, error) {
	var endpointUrl = fmt.Sprintf("https://%s.api.cognitive.microsoft.com/luis/v2.0/apps/%s?subscription-key=%s&verbose=false&q=%s", c.region, c.appID, c.endpointKey, url.QueryEscape(utterance))

	response, err := http.Get(endpointUrl)
	if response.StatusCode != 200 {
		errMsg := ""
		if err != nil {
			errMsg = err.Error()
		}
		return PredictionResult{}, errors.New("None 200 status code " + errMsg)
	}

	if err != nil {
		return PredictionResult{}, err
	}
	resp := ResponseStruct{}
	err = json.NewDecoder(response.Body).Decode(&resp)

	if err != nil {
		return PredictionResult{}, err
	}

	if resp.TopScoringIntent.Score <= 0.5 || resp.TopScoringIntent.Name == "None" {
		return PredictionResult{Result: "unknown"}, nil
	}
	return PredictionResult{Result: resp.TopScoringIntent.Name}, nil
}
