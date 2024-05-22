package api

import "encoding/json"

const EndpointSearchGacha ApiEndpoint = "/searchGacha"

type RequestSearchGacha struct {
	DefaultServers []Server `json:"default_servers"`
	GachaId        int      `json:"gachaId"`
	UseEasyBG      bool     `json:"useEasyBG"`
	Compress       bool     `json:"compress,omitempty"`
}

func searchGachaApi(api *Api, defaultServers []Server, gachaId int, useEasyBG, compress bool) ([]*ApiResponse, error) {
	request := RequestSearchGacha{
		DefaultServers: defaultServers,
		GachaId:        gachaId,
		UseEasyBG:      useEasyBG,
		Compress:       compress,
	}

	response, err := api.post(EndpointSearchGacha, request)
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

func (api *Api) SearchGacha(defaultServers []Server, gachaId int) ([]*ApiResponse, error) {
	return searchGachaApi(api, defaultServers, gachaId, api.config.UseEasyBG, api.config.Compress)
}
