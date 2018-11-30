package net

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HttpMethod string

const (
	HttpMethodPost   HttpMethod = "POST"
	HttpMethodGet    HttpMethod = "GET"
	HttpMethodDelete HttpMethod = "DELETE"
	HttpMethodPut    HttpMethod = "PUT"
)

func HttpHeader() http.Header {
	return http.Header{}
}

func (httpMethod *HttpMethod) String() string {
	return string(*httpMethod)
}

func HttpResponseJsonTo(res *http.Response, to interface{}, _close bool) error {
	if _close {
		defer res.Body.Close()
	} else {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		defer func() {
			res.Body.Close()
			res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}()
	}
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(to)
}

func HttpBasicAuth(username string, password string) string {
	auth := username + ":" + password
	return fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(auth)))
}

func HttpSendPostForm(url string, data url.Values) (*http.Response, error) {
	header := HttpHeader()
	header.Add("Content-Type", "application/x-www-form-urlencoded")
	return HttpSendPost(
		url,
		header,
		strings.NewReader(data.Encode()),
	)
}

func HttpSendPost(url string, header http.Header, body io.Reader) (*http.Response, error) {
	return HttpSend(
		nil,
		url,
		HttpMethodPost,
		header,
		body,
	)
}

func HttpSend(req *http.Request, url string, method HttpMethod, header http.Header, body io.Reader) (*http.Response, error) {
	var err error
	c := &http.Client{}

	if req == nil {
		if req, err = http.NewRequest(method.String(), url, body); err != nil {
			return nil, err
		}
	}

	req.Header = header

	return c.Do(req)
}

type InterceptorResponse struct {
	Data []byte
	Code int
	w    http.ResponseWriter
}

func (w *InterceptorResponse) ResponseWriter(rw http.ResponseWriter) {
	w.w = rw
}

func (w *InterceptorResponse) Header() http.Header {
	return w.w.Header()
}

func (w *InterceptorResponse) Write(b []byte) (int, error) {
	w.Data = b
	return w.w.Write(b)
}

func (w *InterceptorResponse) WriteHeader(i int) {
	w.Code = i
	w.w.WriteHeader(i)
}
