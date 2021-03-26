package main

type AlexaRequest struct {
	Version string `json:"version"`
	Request `json:"request"`
}

type Request struct {
	Intent `json:"intent"`
}

type Intent struct {
	Name               string `json:"name"`
	ConfirmationStatus string `json:"confirmationStatus"`
}
