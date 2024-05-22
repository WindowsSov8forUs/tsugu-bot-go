package api

import "encoding/json"

const EndpointYcx ApiEndpoint = "/ycx"

type RequestYcx struct {
	Server   Server `json:"server"`
	Tier     int    `json:"tier"`
	EventId  int    `json:"eventId,omitempty"`
	Compress bool   `json:"compress,omitempty"`
}

func ycxApi(api *Api, server Server, tier int, eventId int, compress bool) ([]*ApiResponse, error) {
	request := RequestYcx{
		Server:   server,
		Tier:     tier,
		EventId:  eventId,
		Compress: compress,
	}

	response, err := api.post(EndpointYcx, request)
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

func (api *Api) Ycx(server Server, tier int, eventId int, compress bool) ([]*ApiResponse, error) {
	return ycxApi(api, server, tier, eventId, compress)
}
