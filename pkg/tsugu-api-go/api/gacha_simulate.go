package api

import "encoding/json"

const EndpointGachaSimulate ApiEndpoint = "/gachaSimulate"

type RequestGachaSimulate struct {
	ServerMode Server `json:"server_mode"`
	Times      int    `json:"times,omitempty"`
	Compress   bool   `json:"compress,omitempty"`
	GachaId    int    `json:"gachaId,omitempty"`
}

func gachaSimulateApi(api *Api, serverMode Server, times int, compress bool, gachaId int) ([]*ApiResponse, error) {
	request := RequestGachaSimulate{
		ServerMode: serverMode,
		Times:      times,
		Compress:   compress,
		GachaId:    gachaId,
	}

	response, err := api.post(EndpointGachaSimulate, request)
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

func (api *Api) GachaSimulate(serverMode Server, times int, gachaId int) ([]*ApiResponse, error) {
	return gachaSimulateApi(api, serverMode, times, api.config.Compress, gachaId)
}
