package cmd

import (
	"errors"
	"fmt"
	"net/http"
)

type MockHTTP struct {
	Requests  []http.Request
	Responses map[string][]Response
}
type Response struct {
	Response      http.Response
	ResponseError error
}

func (c *MockHTTP) AddResponse(path string, resp Response) {
	c.Responses[path] = append(c.Responses[path], resp)
}
func (c *MockHTTP) pop(path string) (Response, error) {

	arr, ok := c.Responses[path]
	if !ok {
		return Response{}, errors.New(fmt.Sprintf("URL %s not found in mapping", path))
	}
	if len(arr) == 0 {
		return Response{}, errors.New(fmt.Sprintf("URL %s has no more responses", path))
	}

	response, remainingResponseList := arr[0], arr[1:]
	c.Responses[path] = remainingResponseList

	return response, nil

}
func (c *MockHTTP) Do(req *http.Request) (*http.Response, error) {
	c.Requests = append(c.Requests, *req)
	response, err := c.pop(req.URL.Path)

	if err != nil {
		return nil, err
	}
	return &response.Response, response.ResponseError
}

func (c *MockHTTP) CountRequests(path string) int {
	count := 0
	for _, req := range c.Requests {
		if req.URL.Path == path {
			count += 1
		}
	}
	return count
}
