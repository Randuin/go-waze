package waze

import (
  "fmt"
  "net/url"
  "net/http"
  "strings"
  "encoding/json"
  "io/ioutil"
  "bytes"
)

type WazeRequest struct {
  From *Coord
  To *Coord
}

func (request WazeRequest) GetResponse() *WazeResponse {
  resp, _ := http.Get(request.BuildURL())
  response := WazeResponseConstructor(resp)
  return response
}

func (request WazeRequest) BuildURL() string {
  str := fmt.Sprintf("https://www.waze.com/RoutingManager/routingRequest?from=%s&to=%s&at=0&returnJSON=true&returnGeometries=true&returnInstructions=true&timeout=60000&nPaths=3",
    request.From.MakeString(),
    request.To.MakeString())

  return str
}


type WazeResponse struct {
  Alternatives []Alternative `json:"alternatives"`
}

func WazeResponseConstructor(response *http.Response) *WazeResponse {
  str, _ := ioutil.ReadAll(response.Body)
  str = bytes.Replace(str, []byte("NaN"), []byte("null"), -1)
  resp := &WazeResponse{}

  json.Unmarshal(str, &resp)
  return resp
}

type Alternative struct {
  Response struct {
    RouteName string `json:"routeName"`
    Results []Result `json:"results"`
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
  X float64
  Y float64
}

func (c Coord) MakeString() string {
  str := fmt.Sprintf("x:%f+y:%f", c.X, c.Y)
  str = url.QueryEscape(str)
  str = strings.Replace(str, "%2B", "+", 1)
  return str
}
