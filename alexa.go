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

func (r *Response) Say(msg string) { r.Text = msg }

func newResponse() AlexaResponse {
	var resp AlexaResponse
	resp.Version = "1.0"
	resp.Type = "PlainText"
	return resp
}
