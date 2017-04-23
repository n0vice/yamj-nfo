package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SearchContract struct {
	Items []Item `json:"items"`
}

type Item struct {
	Link string `json:"link"`
}

func searchMovie(name string) (*SearchContract, error) {
	client := &http.Client{}
	if *proxy != "" {
		println("creating proxy. ", proxy)
		// create http client with proxy
		proxyURL, err := url.Parse(*proxy)
		if err != nil {
			return nil, fmt.Errorf("failed to create proxy. %v", err.Error())
		}
		client.Transport = &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	}

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%v&cx=%v&q=%v", *googleAPIkey, *googleAPICX, url.QueryEscape(name)), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create search request. %v", err.Error())
	}
	httpResponse, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to get request from Google search. %v", err.Error())
	}
	println(" ", httpResponse.Status)
	defer httpResponse.Body.Close()
	responseContent, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body. %v", err.Error())
	}
	responseMap := &SearchContract{}
	err = json.Unmarshal(responseContent, responseMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body. %v", err.Error())
	}
	return responseMap, nil
}
