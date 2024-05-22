package api

import "encoding/json"

const EndpointRoomList ApiEndpoint = "/roomList"

type RequestRoomList struct {
	RoomList []*Room `json:"roomList"`
	Compress bool    `json:"compress,omitempty"`
}

func roomListApi(api *Api, roomList []*Room, compress bool) ([]*ApiResponse, error) {
	request := RequestRoomList{
		RoomList: roomList,
		Compress: compress,
	}

	response, err := api.post(EndpointRoomList, request)
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

func (api *Api) RoomList(roomList []*Room) ([]*ApiResponse, error) {
	return roomListApi(api, roomList, api.config.Compress)
}
