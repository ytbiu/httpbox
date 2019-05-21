package httpbox

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type Requester interface {
	URL() string
	Params() map[string]string
}

type Request struct {
	url    string
	params map[string]string
}

func (r *Request) URL() string {
	return r.url
}

func (r *Request) Params() map[string]string {
	return r.params
}

func DefaultRequester(url string, params ...map[string]string) Requester {
	r := &Request{url: url}

	if len(params) > 0 {
		r.params = params[0]
	}

	return r
}

func Get(testReq Requester, newHandler http.Handler, resp interface{}) (int, error) {

	params := testReq.Params()
	url := testReq.URL()

	firstLoop := true
	for k, v := range params {
		if firstLoop {
			url += "?"
			firstLoop = false
		}
		url += k + "=" + v + "&"
	}

	req := httptest.NewRequest(http.MethodGet, url[:len(url)-1], nil)
	body, respCode := sendReq(req, newHandler)

	return respCode, json.Unmarshal(body, resp)
}

func Post(testReq Requester, newHandler http.Handler, resp interface{}) (int, error) {

	params := testReq.Params()
	url := testReq.URL()

	postForm, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(postForm))
	body, respCode := sendReq(req, newHandler)

	return respCode, json.Unmarshal(body, resp)
}

func sendReq(req *http.Request, handler http.Handler) ([]byte, int) {

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	respBody, _ := ioutil.ReadAll(result.Body)
	return respBody, result.StatusCode
}
