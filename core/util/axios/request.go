package axios

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/dalton02/licor/httpkit"
)

type Axios struct {
	BaseUrl string
	Headers map[string]string
	Client  *http.Client
}

func (axios *Axios) NewRequest(method string, path string, data any) (*http.Request, error) {
	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, axios.BaseUrl+path, body)
	if err != nil {
		return req, err
	}
	for key, value := range axios.Headers {
		req.Header.Set(key, value)
	}

	return req, err
}

func (axios *Axios) Do(req *http.Request) (*http.Response, error) {

	resp, err := axios.Client.Do(req)
	return resp, err
}
func (axios *Axios) ParseResponseBody(res *http.Response) (map[string]interface{}, error) {
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (axios *Axios) ParseHttpKitResponse(res *http.Response) (httpkit.HttpMessage, bool) {
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)

	var result httpkit.HttpMessage
	json.Unmarshal(bodyBytes, &result)

	return result, false
}

func (axios *Axios) SetHeaders(headers map[string]string) {
	axios.Headers = headers
}
func (axios *Axios) SetBaseUrl(baseUrl string) {
	axios.BaseUrl = baseUrl
}
