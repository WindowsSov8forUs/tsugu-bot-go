package api

import "encoding/json"

const EndpointSearchEvent ApiEndpoint = "/searchEvent"

type RequestSearchEvent struct {
	DefaultServers []Server `json:"default_servers"`
	Text           string   `json:"text"`
	UseEasyBG      bool     `json:"useEasyBG"`
	Compress       bool     `json:"compress,omitempty"`
}

func searchEventApi(api *Api, defaultServers []Server, text string, useEasyBG, compress bool) ([]*ApiResponse, error) {
	request := RequestSearchEvent{
		DefaultServers: defaultServers,
		Text:           text,
		UseEasyBG:      useEasyBG,
		Compress:       compress,
	}

	response, err := api.post(EndpointSearchEvent, request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var apiResponse []*ApiResponse
	err = json.NewDecoder(response.Body).Decode(&apiResponse)
	if err != nil {
		return nil, err
	}

	return apiResponse, nil
}

func (api *Api) SearchEvent(defaultServers []Server, text string) ([]*ApiResponse, error) {
	return searchEventApi(api, defaultServers, text, api.config.UseEasyBG, api.config.Compress)
}
