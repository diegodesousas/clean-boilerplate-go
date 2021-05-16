package mock

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gojektech/heimdall/v6"
)

type mockParams struct {
	response *http.Response
	err      error
}

type HttpClient struct {
	counterGet    int
	queueGet      map[int]mockParams
	counterPost   int
	queuePost     map[int]mockParams
	counterPut    int
	queuePut      map[int]mockParams
	counterPatch  int
	queuePatch    map[int]mockParams
	counterDelete int
	queueDelete   map[int]mockParams
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		queueGet:  map[int]mockParams{},
		queuePost: map[int]mockParams{},
		queuePut:  map[int]mockParams{},
	}
}

func (h *HttpClient) AddPostResponse(res *http.Response, err error) *HttpClient {
	h.queuePost[len(h.queuePost)] = mockParams{res, err}

	return h
}

func (h *HttpClient) AddGetResponse(res *http.Response, err error) *HttpClient {
	h.queueGet[len(h.queueGet)] = mockParams{res, err}

	return h
}

func (h *HttpClient) AddPutResponse(res *http.Response, err error) *HttpClient {
	h.queuePut[len(h.queuePut)] = mockParams{res, err}

	return h
}

func (h HttpClient) Get(url string, headers http.Header) (*http.Response, error) {
	m, ok := h.queueGet[h.counterGet]
	if ok {
		h.counterGet += 1
		return m.response, m.err
	}

	panic("httpclient-mock: get queue is empty")
}

func (h *HttpClient) Post(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	m, ok := h.queuePost[h.counterPost]
	if ok {
		h.counterPost += 1
		return m.response, m.err
	}

	panic("httpclient-mock: post queue is empty")
}

func (h HttpClient) Put(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	m, ok := h.queuePut[h.counterPut]
	if ok {
		h.counterPut += 1
		return m.response, m.err
	}

	panic("httpclient-mock: put queue is empty")
}

func (h HttpClient) Patch(url string, body io.Reader, headers http.Header) (*http.Response, error) {
	panic("implement me")
}

func (h HttpClient) Delete(url string, headers http.Header) (*http.Response, error) {
	panic("implement me")
}

func (h HttpClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("do method not implemented")
}

func (h HttpClient) AddPlugin(p heimdall.Plugin) {}

func MarshalJson(data interface{}) io.ReadCloser {
	dataBytes, _ := json.Marshal(data)

	return ioutil.NopCloser(bytes.NewReader(dataBytes))
}

func MarshalString(data string) io.ReadCloser {
	return ioutil.NopCloser(strings.NewReader(data))
}
