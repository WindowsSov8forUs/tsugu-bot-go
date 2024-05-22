package api

import "encoding/json"

const EndpointSongChart ApiEndpoint = "/songChart"

type RequestSongChart struct {
	DefaultServers []Server `json:"default_servers"`
	SongId         int      `json:"songId"`
	DifficultyText string   `json:"difficultyText"`
	Compress       bool     `json:"compress,omitempty"`
}

func songChartApi(api *Api, defaultServers []Server, songId int, difficultyText string, compress bool) ([]*ApiResponse, error) {
	request := RequestSongChart{
		DefaultServers: defaultServers,
		SongId:         songId,
		DifficultyText: difficultyText,
		Compress:       compress,
	}

	response, err := api.post(EndpointSongChart, request)
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

func (api *Api) SongChart(defaultServers []Server, songId int, difficultyText string) ([]*ApiResponse, error) {
	return songChartApi(api, defaultServers, songId, difficultyText, api.config.Compress)
}
