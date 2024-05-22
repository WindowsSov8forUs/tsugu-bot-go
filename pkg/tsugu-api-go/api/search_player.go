package api

import "encoding/json"

const EndpointSearchPlayer ApiEndpoint = "/searchPlayer"

type RequestSearchPlayer struct {
	PlayerId  int    `json:"playerId"`
	Server    Server `json:"server"`
	UseEasyBG bool   `json:"useEasyBG"`
	Compress  bool   `json:"compress,omitempty"`
}

func searchPlayerApi(api *Api, playerId int, server Server, useEasyBG, compress bool) ([]*ApiResponse, error) {
	request := RequestSearchPlayer{
		PlayerId:  playerId,
		Server:    server,
		UseEasyBG: useEasyBG,
		Compress:  compress,
	}

	response, err := api.post(EndpointSearchPlayer, request)
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

func (api *Api) SearchPlayer(playerId int, server Server) ([]*ApiResponse, error) {
	return searchPlayerApi(api, playerId, server, api.config.UseEasyBG, api.config.Compress)
}
