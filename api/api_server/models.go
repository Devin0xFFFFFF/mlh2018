package api_server

type ResponseStruct struct {
	Query            string         `json:"query"`
	TopScoringIntent IntentResult   `json:"topScoringIntent"`
	Entities         []EntityResult `json:"entities"`
}

type IntentResult struct {
	Name  string  `json:"intent"`
	Score float64 `json:"score"`
}
type EntityResult struct {
	Name string `json:"name"`
}

type PredictionResult struct {
	Result string `json:"result"`
}

type IntentResponse struct {
	TopScoringIntent IntentResult   `json:"intent"`
	Entities         []EntityResult `json:"entities"`
}

type RecognitionResponse struct {
	Status      string `json:"RecognitionStatus"`
	DisplayText string `json:"DisplayText"`
	Offset      int    `json:"Offset"`
	Duration    int    `json:"Duration"`
}
