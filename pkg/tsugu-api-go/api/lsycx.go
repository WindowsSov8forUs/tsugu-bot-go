package api

import "encoding/json"

const EndpointLsycx ApiEndpoint = "/lsycx"

type RequestLsycx struct {
	Server   Server `json:"server"`
	Tier     int    `json:"tier"`
	EventId  int    `json:"eventId,omitempty"`
	Compress bool   `json:"compress,omitempty"`
}

func lsycxApi(api *Api, server Server, tier int, eventId int, compress bool) ([]*ApiResponse, error) {
	request := RequestLsycx{
		Server:   server,
		Tier:     tier,
		EventId:  eventId,
		Compress: compress,
	}

	response, err := api.post(EndpointLsycx, request)
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

func (api *Api) Lsycx(server Server, tier int, eventId int) ([]*ApiResponse, error) {
	return lsycxApi(api, server, tier, eventId, api.config.Compress)
}
