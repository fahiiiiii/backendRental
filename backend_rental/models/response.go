package models

type PropertyResponse struct {
    Properties []Property `json:"properties"`
    Status     string    `json:"status"`
    Message    string    `json:"message"`
}

type PropertyResult struct {
    City       CityKey
    Properties []Property
    Err        error
}
