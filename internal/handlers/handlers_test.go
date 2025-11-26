package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	client = http.Client{Timeout: time.Second * 10}
)

type Case struct {
	Name     string
	URL      string
	Status   int
	Body     interface{}
	Expected interface{}
}

func performRequest(client http.Client, url string, body interface{}) (int, interface{}, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return 0, nil, err
	}

	reqBody := bytes.NewReader(data)
	req, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("request error: %v", err)
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("error readall: %v", err)
	}

	var result interface{}
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return 0, nil, fmt.Errorf("cant unpack json: %v", err)
	}

	return resp.StatusCode, result, nil
}

func TestAPI(t *testing.T) {
	cases := []Case{
		Case{
			Name:     "check one link",
			URL:      "/links",
			Status:   http.StatusOK,
			Body:     map[string]interface{}{"links": []string{"google.com"}},
			Expected: map[string]interface{}{"links": map[string]interface{}{"google.com": "available"}, "links_num": 0},
		},
		Case{
			Name:   "check two links",
			URL:    "/links",
			Status: http.StatusOK,
			Body:   map[string]interface{}{"links": []string{"not_available.com", "yandex.ru"}},
			Expected: map[string]interface{}{
				"links": map[string]interface{}{"not_available.com": "not available", "yandex.ru": "available"}, "links_num": 1,
			},
		},
		Case{
			Name:   "check double links",
			URL:    "/links",
			Status: http.StatusOK,
			Body:   map[string]interface{}{"links": []string{"yandex.ru", "yandex.ru"}},
			Expected: map[string]interface{}{
				"links": map[string]interface{}{"yandex.ru": "available"}, "links_num": 2,
			},
		},
	}

	api := NewAPI()
	server := httptest.NewServer(NewMux(api))

	runCases(t, server, cases)
}

func runCases(t *testing.T, ts *httptest.Server, cases []Case) {
	for idx, item := range cases {
		caseName := fmt.Sprintf("case %d: [%s]", idx, item.Name)

		code, resp, err := performRequest(client, ts.URL+item.URL, item.Body)
		if err != nil {
			t.Fatalf("[%s] %s", caseName, err)
		}

		assert.Equal(t, item.Status, code)

		jsonExpected, err := json.Marshal(item.Expected)
		if err != nil {
			t.Fatalf("[%s] cant marshal expected json: %v", caseName, err)
		}

		var interfaceExpected interface{}
		err = json.Unmarshal(jsonExpected, &interfaceExpected)
		if err != nil {
			t.Fatalf("[%s] cant unmarshal expected json: %v", caseName, err)
		}

		if !reflect.DeepEqual(resp, interfaceExpected) {
			t.Errorf("[%s] results not match\nGot : %#v\nExpected: %#v", caseName, resp, interfaceExpected)
		}
	}

}
