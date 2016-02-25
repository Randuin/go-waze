package waze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	From *Coord
	To   *Coord
}

func (request Request) GetResponse() *Response {
	resp, _ := http.Get(request.BuildURL())
	response := ResponseConstructor(resp)
	return response
}

func (request Request) BuildURL() string {
	str := fmt.Sprintf("https://www.waze.com/RoutingManager/routingRequest?from=%s&to=%s&at=0&returnJSON=true&returnGeometries=true&returnInstructions=true&timeout=60000&nPaths=3",
		request.From.MakeString(),
		request.To.MakeString())

	return str
}

type Response struct {
	Alternatives []Alternative `json:"alternatives"`
}

func ResponseConstructor(response *http.Response) *Response {
	str, _ := ioutil.ReadAll(response.Body)
	str = bytes.Replace(str, []byte("NaN"), []byte("null"), -1)
	resp := &Response{}

	json.Unmarshal(str, &resp)
	return resp
}

type Alternative struct {
	Response struct {
		RouteName string   `json:"routeName"`
		Results   []Result `json:"results"`
	}
}

func (a Alternative) TotalCrossTime() int {
	sum := 0
	for _, result := range a.Response.Results {
		sum += result.CrossTime
	}
	return sum
}

type Result struct {
	CrossTime int `json:"crossTime"`
}

type Coord struct {
	Lat float64
	Lng float64
}

func (c Coord) MakeString() string {
	str := fmt.Sprintf("x:%f+y:%f", c.Lng, c.Lat)
	str = url.QueryEscape(str)
	str = strings.Replace(str, "%2B", "+", 1)
	return str
}
