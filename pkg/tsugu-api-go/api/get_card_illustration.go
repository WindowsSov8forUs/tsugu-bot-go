package api

import "encoding/json"

const EndpointGetCardIllustration ApiEndpoint = "/getCardIllustration"

type RequestGetCardIllustration struct {
	CardId int `json:"cardId"`
}

func getCardIllustrationApi(api *Api, cardId int) ([]*ApiResponse, error) {
	request := RequestGetCardIllustration{
		CardId: cardId,
	}

	response, err := api.post(EndpointGetCardIllustration, request)
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

func (api *Api) GetCardIllustration(cardId int) ([]*ApiResponse, error) {
	return getCardIllustrationApi(api, cardId)
}
