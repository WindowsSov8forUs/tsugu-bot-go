package api

import "encoding/json"

const EndpointEventStage ApiEndpoint = "/eventStage"

type RequestEventStage struct {
	Server   Server `json:"server"`
	EventId  int    `json:"eventId,omitempty"`
	Meta     bool   `json:"meta,omitempty"`
	Compress bool   `json:"compress,omitempty"`
}

func eventStageApi(api *Api, server Server, eventId int, meta, compress bool) ([]*ApiResponse, error) {
	request := RequestEventStage{
		Server:   server,
		EventId:  eventId,
		Meta:     meta,
		Compress: compress,
	}

	response, err := api.post(EndpointEventStage, request)
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

func (api *Api) EventStage(server Server, eventId int, meta bool) ([]*ApiResponse, error) {
	return eventStageApi(api, server, eventId, meta, api.config.Compress)
}
