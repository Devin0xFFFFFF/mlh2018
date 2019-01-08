package api_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	appID       string
	ocpKey      string
	endpointKey string
	region      string
}

func NewClient(appID string, endpointKey string, region string, ocpKey string) *Client {
	return &Client{appID: appID, endpointKey: endpointKey, region: region, ocpKey: ocpKey}
}
func (c *Client) PredictFromVoice(voiceFile io.Reader) (PredictionResult, error) {

	var endpointUrl = fmt.Sprintf("https://%s.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=en-US", c.region)
	request, err := http.NewRequest("POST", endpointUrl, voiceFile)
	request.Header.Set("Ocp-Apim-Subscription-Key", c.ocpKey)
	request.Header.Set("Content-Type", "audio/wav; codec=audio/pcm; samplerate=16000")
	client := &http.Client{}
	response, err := client.Do(request)
	//response, err := http.Post(endpointUrl, "audio/wav; codec=audio/pcm; samplerate=16000", voiceFile)
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
	resp := RecognitionResponse{}
	err = json.NewDecoder(response.Body).Decode(&resp)
	fmt.Println(resp)
	if err != nil {
		return PredictionResult{}, err
	}
	if resp.Status == "Success" {
		// if word recognition is successful, call predict
		return c.Predict(resp.DisplayText)
	}

	return PredictionResult{}, errors.New("voice recognition failed")

}

func (c *Client) Predict(utterance string) (PredictionResult, error) {
	var endpointUrl = fmt.Sprintf("https://%s.api.cognitive.microsoft.com/luis/v2.0/apps/%s?subscription-key=%s&verbose=false&q=%s", c.region, c.appID, c.endpointKey, url.QueryEscape(utterance))
	fmt.Println("predicting utterance of ", utterance)
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
