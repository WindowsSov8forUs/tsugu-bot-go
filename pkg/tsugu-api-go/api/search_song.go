package api

import "encoding/json"

const EndpointSearchSong ApiEndpoint = "/searchSong"

type RequestSearchSong struct {
	DefaultServers []Server `json:"default_servers"`
	Text           string   `json:"text"`
	Compress       bool     `json:"compress,omitempty"`
}

func searchSongApi(api *Api, defaultServers []Server, text string, compress bool) ([]*ApiResponse, error) {
	request := RequestSearchSong{
		DefaultServers: defaultServers,
		Text:           text,
		Compress:       compress,
	}

	response, err := api.post(EndpointSearchSong, request)
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

func (api *Api) SearchSong(defaultServers []Server, text string) ([]*ApiResponse, error) {
	return searchSongApi(api, defaultServers, text, api.config.Compress)
}
