package api

import "encoding/json"

const EndpointSongMeta ApiEndpoint = "/songMeta"

type RequestSongMeta struct {
	DefaultServers []Server `json:"default_servers"`
	Server         Server   `json:"server"`
	Compress       bool     `json:"compress,omitempty"`
}

func songMetaApi(api *Api, defaultServers []Server, server Server, compress bool) ([]*ApiResponse, error) {
	request := RequestSongMeta{
		DefaultServers: defaultServers,
		Server:         server,
		Compress:       compress,
	}

	response, err := api.post(EndpointSongMeta, request)
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

func (api *Api) SongMeta(defaultServers []Server, server Server) ([]*ApiResponse, error) {
	return songMetaApi(api, defaultServers, server, api.config.Compress)
}
