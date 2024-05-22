package api

import "encoding/json"

const EndpointStationSubmitRoomNumber ApiEndpoint = "/station/submitRoomNumber"

type RequestStationSubmitRoomNumber struct {
	Number              int    `json:"number"`
	RawMessage          string `json:"rawMessage"`
	Platform            string `json:"platform"`
	UserId              string `json:"user_id"`
	UserName            string `json:"userName"`
	Time                int64  `json:"time"`
	BandoriStationToken string `json:"bandoriStationToken,omitempty"`
}

type ResponseStationSubmitRoomNumber struct {
	Status ApiResponseStatus `json:"status"`
	Data   string            `json:"data"`
}

func stationSubmitRoomNumberApi(api *RemoteDBApi, number int, rawMessage string, platform string, userId string, userName string, time int64, bandoriStationToken string) (*ResponseStationSubmitRoomNumber, *ResponseFailed, error) {
	request := RequestStationSubmitRoomNumber{
		Number:              number,
		RawMessage:          rawMessage,
		Platform:            platform,
		UserId:              userId,
		UserName:            userName,
		Time:                time,
		BandoriStationToken: bandoriStationToken,
	}

	response, status, err := api.post(EndpointStationSubmitRoomNumber, request)
	if err != nil {
		return nil, nil, err
	}
	if status {
		var apiResponse *ResponseStationSubmitRoomNumber
		err = json.NewDecoder(response.Body).Decode(&apiResponse)
		if err != nil {
			return nil, nil, err
		}
		return apiResponse, nil, nil
	} else {
		var apiResponse *ResponseFailed
		err = json.NewDecoder(response.Body).Decode(&apiResponse)
		if err != nil {
			return nil, nil, err
		}
		return nil, apiResponse, nil
	}
}

func (api *Api) StationSubmitRoomNumber(number int, rawMessage string, platform string, userId string, userName string, time int64, bandoriStationToken string) (*ResponseStationSubmitRoomNumber, *ResponseFailed, error) {
	return stationSubmitRoomNumberApi(api.remoteDB, number, rawMessage, platform, userId, userName, time, bandoriStationToken)
}

const EndpointStationQueryAllRooms ApiEndpoint = "/station/queryAllRooms"

type ResponseStationQueryAllRooms struct {
	Status ApiResponseStatus `json:"status"`
	Data   []*Room           `json:"data"`
}

func stationQueryAllRoomsApi(api *RemoteDBApi) (*ResponseStationQueryAllRooms, *ResponseFailed, error) {
	response, status, err := api.get(EndpointStationQueryAllRooms, nil)
	if err != nil {
		return nil, nil, err
	}
	if status {
		var apiResponse *ResponseStationQueryAllRooms
		err = json.NewDecoder(response.Body).Decode(&apiResponse)
		if err != nil {
			return nil, nil, err
		}
		return apiResponse, nil, nil
	} else {
		var apiResponse *ResponseFailed
		err = json.NewDecoder(response.Body).Decode(&apiResponse)
		if err != nil {
			return nil, nil, err
		}
		return nil, apiResponse, nil
	}
}

func (api *Api) StationQueryAllRooms() (*ResponseStationQueryAllRooms, *ResponseFailed, error) {
	return stationQueryAllRoomsApi(api.remoteDB)
}
