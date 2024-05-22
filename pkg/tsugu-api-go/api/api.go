package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type Api struct {
	Url      string
	Proxy    string
	Timeout  int
	remoteDB *RemoteDBApi
	config   *Config
}

func (api *Api) post(endpoint ApiEndpoint, data interface{}) (*http.Response, error) {
	var transport *http.Transport
	if api.Proxy != "" {
		proxyUrl, _ := url.Parse(api.Proxy)
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	} else {
		transport = &http.Transport{}
	}

	client := &http.Client{Transport: transport}
	if api.Timeout > 0 {
		client.Timeout = time.Duration(api.Timeout) * time.Second
	}

	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	requestUrl := api.Url + string(endpoint)
	response, err := client.Post(requestUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusBadRequest {
		return response, nil
	} else {
		return nil, &ErrBadStatus{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
	}
}

type RemoteDBApi struct {
	Url     string
	Proxy   string
	Timeout int
}

func (api *RemoteDBApi) get(endpoint ApiEndpoint, params *map[string]string) (*http.Response, bool, error) {
	var transport *http.Transport
	if api.Proxy != "" {
		proxyUrl, _ := url.Parse(api.Proxy)
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	} else {
		transport = &http.Transport{}
	}

	client := &http.Client{Transport: transport}
	if api.Timeout > 0 {
		client.Timeout = time.Duration(api.Timeout) * time.Second
	}

	var requestUrl string
	if params != nil {
		urlParams := url.Values{}
		for key, value := range *params {
			urlParams.Add(key, value)
		}
		requestUrl = api.Url + string(endpoint) + "?" + urlParams.Encode()
	} else {
		requestUrl = api.Url + string(endpoint)
	}

	response, err := client.Get(requestUrl)
	if err != nil {
		return nil, false, err
	}
	if response.StatusCode == http.StatusOK {
		return response, true, nil
	} else if response.StatusCode == http.StatusBadRequest {
		return response, false, nil
	} else {
		return nil, false, &ErrBadStatus{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
	}
}

func (api *RemoteDBApi) post(endpoint ApiEndpoint, data interface{}) (*http.Response, bool, error) {
	var transport *http.Transport
	if api.Proxy != "" {
		proxyUrl, _ := url.Parse(api.Proxy)
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	} else {
		transport = &http.Transport{}
	}

	client := &http.Client{Transport: transport}
	if api.Timeout > 0 {
		client.Timeout = time.Duration(api.Timeout) * time.Second
	}

	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, false, err
	}
	requestUrl := api.Url + string(endpoint)
	response, err := client.Post(requestUrl, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, false, err
	}
	if response.StatusCode == http.StatusOK {
		return response, true, nil
	} else if response.StatusCode == http.StatusBadRequest {
		return response, false, nil
	} else {
		return nil, false, &ErrBadStatus{
			StatusCode: response.StatusCode,
			Message:    response.Status,
		}
	}
}

func NewApi(config *Config) *Api {
	remoteDBApi := &RemoteDBApi{
		Url:     config.DatabaseBackendUrl,
		Timeout: config.Timeout,
	}
	if config.DatabaseBackendProxy {
		remoteDBApi.Proxy = config.Proxy
	}
	api := &Api{
		Url:      config.BackendUrl,
		Timeout:  config.Timeout,
		remoteDB: remoteDBApi,
		config:   config,
	}
	if config.BackendProxy {
		api.Proxy = config.Proxy
	}
	return api
}
