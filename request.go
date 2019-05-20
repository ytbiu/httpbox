package httpbox

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func Get(url string, params map[string]string, newHandler func() http.Handler) ([]byte, int) {

	firstLoop := true
	for k, v := range params {
		if firstLoop {
			url += "?"
			firstLoop = false
		}
		url += k + "=" + v + "&"
	}

	req := httptest.NewRequest(http.MethodGet, url[:len(url)-1], nil)

	return sendReq(req, newHandler)
}

func Post(url string, params map[string]string, newHandler func() http.Handler) ([]byte, int) {

	body, _ := json.Marshal(params)
	req := httptest.NewRequest(http.MethodPost, url, bytes.NewReader(body))

	return sendReq(req, newHandler)
}

func sendReq(req *http.Request, handler func() http.Handler) ([]byte, int) {

	w := httptest.NewRecorder()
	handler().ServeHTTP(w, req)

	result := w.Result()
	defer result.Body.Close()

	respBody, _ := ioutil.ReadAll(result.Body)
	return respBody, result.StatusCode
}
