package api

import "encoding/json"

const EndpointYcxAll ApiEndpoint = "/ycxAll"

type RequestYcxAll struct {
	Server   Server `json:"server"`
	EventId  int    `json:"eventId,omitempty"`
	Compress bool   `json:"compress,omitempty"`
}

func ycxAllApi(api *Api, server Server, eventId int, compress bool) ([]*ApiResponse, error) {
	request := RequestYcxAll{
		Server:   server,
		EventId:  eventId,
		Compress: compress,
	}

	response, err := api.post(EndpointYcxAll, request)
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

func (api *Api) YcxAll(server Server, eventId int, compress bool) ([]*ApiResponse, error) {
	return ycxAllApi(api, server, eventId, compress)
}
