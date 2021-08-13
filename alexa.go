package main

type AlexaResponse struct {
	Response `json:"response"`
	Version  string `json:"version"`
}

type Response struct {
	OutputSpeech `json:"outputSpeech"`
}

type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

func newResponse(msg string) *AlexaResponse {
	return &AlexaResponse{
		Response: Response{
			OutputSpeech: OutputSpeech{
				Type: "PlainText",
				Text: msg,
			},
		},
		Version: "1.0",
	}
}
