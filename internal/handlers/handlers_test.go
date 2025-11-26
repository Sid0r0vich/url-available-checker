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
)

type Case struct {
	Name     string
	URL      string
	Status   int
	Body     interface{}
	Expected interface{}
}

func TestAPI(t *testing.T) {
	cases := []Case{
		Case{
			Name:     "check /links",
			URL:      "/links",
			Status:   http.StatusOK,
			Body:     map[string]interface{}{"links": []string{"google.com"}},
			Expected: map[string]interface{}{"links": map[string]interface{}{"google.com": "available"}, "links_num": 0},
		},
		Case{
			Name:   "check /links",
			URL:    "/links",
			Status: http.StatusOK,
			Body:   map[string]interface{}{"links": []string{"not_available.com", "google.com"}},
			Expected: map[string]interface{}{
				"links": map[string]interface{}{"not_available.com": "not available", "google.com": "available"}, "links_num": 1,
			},
		},
	}

	r := http.NewServeMux()
	api := NewAPI()

	r.HandleFunc("/links", api.GetLinksHanlder)
	r.HandleFunc("/list", api.MakePDFHandler)

	server := httptest.NewServer(r)

	runCases(t, server, cases)
}

func runCases(t *testing.T, ts *httptest.Server, cases []Case) {
	client := &http.Client{Timeout: time.Second}

	for idx, item := range cases {
		caseName := fmt.Sprintf("case %d: [%s]", idx, item.Name)

		data, err := json.Marshal(item.Body)
		if err != nil {
			t.Fatal(err)
		}

		reqBody := bytes.NewReader(data)
		var errNewReq error
		url := ts.URL + item.URL
		fmt.Printf("URL: %s\n", url)
		req, err := http.NewRequest("POST", url, reqBody)
		if err != nil {
			t.Fatal(errNewReq)
		}
		req.Header.Add("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("[%s] request error: %v", caseName, err)
			continue
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("[%s] error readall: %v", caseName, err)
			continue
		}

		if resp.StatusCode != item.Status {
			t.Fatalf("[%s] expected http status %v, got %v", caseName, item.Status, resp.StatusCode)
			continue
		}

		var result interface{}
		err = json.Unmarshal(body, &result)
		if err != nil {
			t.Fatalf("[%s] cant unpack json: %v", caseName, err)
			continue
		}

		jsonExpected, err := json.Marshal(item.Expected)
		if err != nil {
			t.Fatalf("[%s] cant marshal expected json: %v", caseName, err)
			continue
		}

		var interfaceExpected interface{}
		err = json.Unmarshal(jsonExpected, &interfaceExpected)
		if err != nil {
			t.Fatalf("[%s] cant unmarshal expected json: %v", caseName, err)
		}

		if !reflect.DeepEqual(result, interfaceExpected) {
			t.Fatalf("[%s] results not match\nGot : %#v\nExpected: %#v", caseName, result, interfaceExpected)
			continue
		}
	}

}
