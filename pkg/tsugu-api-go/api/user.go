package api

import "encoding/json"

const EndpointUserGetUserData ApiEndpoint = "/user/getUserData"

type RequestUserGetUserData struct {
	Platform string `json:"platform"`
	UserId   string `json:"user_id"`
}

type ResponseUserGetUserData struct {
	Status ApiResponseStatus `json:"status"`
	Data   *tsuguUser        `json:"data"`
}

func userGetUserDataApi(api *RemoteDBApi, platform string, userId string) (*ResponseUserGetUserData, *ResponseFailed, error) {
	request := RequestUserGetUserData{
		Platform: platform,
		UserId:   userId,
	}

	response, status, err := api.post(EndpointUserGetUserData, request)
	if err != nil {
		return nil, nil, err
	}
	if status {
		var apiResponse *ResponseUserGetUserData
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

func (api *Api) UserGetUserData(platform string, userId string) (*ResponseUserGetUserData, *ResponseFailed, error) {
	return userGetUserDataApi(api.remoteDB, platform, userId)
}

const EndpointUserChangeUserData ApiEndpoint = "/user/changeUserData"

type RequestUserChangeUserData struct {
	Platform string           `json:"platform"`
	UserId   string           `json:"user_id"`
	Update   PartialTsuguUser `json:"update"`
}

type ResponseUserChangeUserData struct {
	Status ApiResponseStatus `json:"status"`
}

func userChangeUserDataApi(api *RemoteDBApi, platform string, userId string, update PartialTsuguUser) (*ResponseUserChangeUserData, *ResponseFailed, error) {
	request := RequestUserChangeUserData{
		Platform: platform,
		UserId:   userId,
		Update:   update,
	}

	response, status, err := api.post(EndpointUserChangeUserData, request)
	if err != nil {
		return nil, nil, err
	}
	if status {
		var apiResponse *ResponseUserChangeUserData
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

func (api *Api) UserChangeUserData(platform string, userId string, update PartialTsuguUser) (*ResponseUserChangeUserData, *ResponseFailed, error) {
	return userChangeUserDataApi(api.remoteDB, platform, userId, update)
}

const EndpointUserBindPlayerRequest ApiEndpoint = "/user/bindPlayerRequest"

type RequestUserBindPlayerRequest struct {
	Platform string `json:"platform"`
	UserId   string `json:"user_id"`
	Server   Server `json:"server"`
	BindType bool   `json:"bindType"`
}

type ResponseUserBindPlayerRequest struct {
	Status ApiResponseStatus `json:"status"`
	Data   struct {
		VerifyCode int `json:"verifyCode"`
	} `json:"data"`
}

func userBindPlayerRequestApi(api *RemoteDBApi, platform string, userId string, server Server, bindType bool) (*ResponseUserBindPlayerRequest, *ResponseFailed, error) {
	request := RequestUserBindPlayerRequest{
		Platform: platform,
		UserId:   userId,
		Server:   server,
		BindType: bindType,
	}

	response, status, err := api.post(EndpointUserBindPlayerRequest, request)
	if err != nil {
		return nil, nil, err
	}
	if status {
		var apiResponse *ResponseUserBindPlayerRequest
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

func (api *Api) UserBindPlayerRequest(platform string, userId string, server Server, bindType bool) (*ResponseUserBindPlayerRequest, *ResponseFailed, error) {
	return userBindPlayerRequestApi(api.remoteDB, platform, userId, server, bindType)
}

const EndpointUserBindPlayerVerification ApiEndpoint = "/user/bindPlayerVerification"

type RequestUserBindPlayerVerification struct {
	Platform string `json:"platform"`
	UserId   string `json:"user_id"`
	Server   Server `json:"server"`
	PlayerId int    `json:"playerId"`
	BindType bool   `json:"bindType"`
}

type ResponseUserBindPlayerVerification struct {
	Status ApiResponseStatus `json:"status"`
	Data   string            `json:"data"`
}

func userBindPlayerVerificationApi(api *RemoteDBApi, platform string, userId string, server Server, playerId int, bindType bool) (*ResponseUserBindPlayerVerification, *ResponseFailed, error) {
	request := RequestUserBindPlayerVerification{
		Platform: platform,
		UserId:   userId,
		Server:   server,
		PlayerId: playerId,
		BindType: bindType,
	}

	response, status, err := api.post(EndpointUserBindPlayerVerification, request)
	if err != nil {
		return nil, nil, err
	}
	if status {
		var apiResponse *ResponseUserBindPlayerVerification
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

func (api *Api) UserBindPlayerVerification(platform string, userId string, server Server, playerId int, bindType bool) (*ResponseUserBindPlayerVerification, *ResponseFailed, error) {
	return userBindPlayerVerificationApi(api.remoteDB, platform, userId, server, playerId, bindType)
}
