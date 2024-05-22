package api

import "encoding/json"

const EndpointSearchCharacter ApiEndpoint = "/searchCharacter"

type RequestSearchCharacter struct {
	DefaultServers []Server `json:"default_servers"`
	Text           string   `json:"text"`
	Compress       bool     `json:"compress,omitempty"`
}

func searchCharacterApi(api *Api, defaultServers []Server, text string, compress bool) ([]*ApiResponse, error) {
	request := RequestSearchCharacter{
		DefaultServers: defaultServers,
		Text:           text,
		Compress:       compress,
	}

	response, err := api.post(EndpointSearchCharacter, request)
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

func (api *Api) SearchCharacter(defaultServers []Server, text string) ([]*ApiResponse, error) {
	return searchCharacterApi(api, defaultServers, text, api.config.Compress)
}
