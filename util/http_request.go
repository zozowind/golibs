package util

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultHTTPClient default http client
var DefaultHTTPClient *http.Client

// DefaultMediaHTTPClient default http client for request media
var DefaultMediaHTTPClient *http.Client

func init() {
	shortTimeClient := *http.DefaultClient
	shortTimeClient.Timeout = time.Second * 30
	DefaultHTTPClient = &shortTimeClient

	longTimeclient := *http.DefaultClient
	longTimeclient.Timeout = time.Second * 60
	DefaultMediaHTTPClient = &longTimeclient
}

//RequestFile upload media file strcut
type RequestFile struct {
	Name string
	Data *bytes.Buffer
}

func (f *RequestFile) Read(p []byte) (n int, err error) {
	return f.Data.Read(p)
}

//HTTPFormPost http request using form post
func HTTPFormPost(httpClient *http.Client, url string, params url.Values, files map[string][]*RequestFile) ([]byte, error) {
	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}
	//create an empty form
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//add string params

	for key, s := range params {
		for _, v := range s {
			_ = bodyWriter.WriteField(key, v)
		}
	}

	//add file to upload
	for key, fs := range files {
		for _, f := range fs {
			//create file field
			fileWriter, err := bodyWriter.CreateFormFile(key, f.Name)
			if nil != err {
				return nil, err
			}
			//copy filedata to form
			_, err = io.Copy(fileWriter, f)
			if err != nil {
				return nil, err
			}
		}
	}

	// get upload content-type like multipart/form-data; boundary=...
	contentType := bodyWriter.FormDataContentType()

	// close bodyWriter now, not in deferr, it will add close tag to body
	bodyWriter.Close()

	response := []byte{}
	// post to server
	resp, err := httpClient.Post(url, contentType, bodyBuf)
	if err != nil {
		return response, err
	}

	response, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return response, err
	}
	return response, nil
}

//HTTPPost x-www-form-urlencoded post
func HTTPPost(httpClient *http.Client, url string, params url.Values) ([]byte, error) {
	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}
	var bodyBuf io.Reader
	if params != nil {
		bodyBuf = strings.NewReader(params.Encode())
	}
	resp, err := httpClient.Post(url, "application/x-www-form-urlencoded", bodyBuf)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, err
}

// HTTPGet http get request
func HTTPGet(httpClient *http.Client, url string, headers map[string]string) ([]byte, error) {
	data := []byte{}
	if httpClient == nil {
		httpClient = DefaultHTTPClient
	}
	req, err := http.NewRequest("GET", url, nil)
	if nil != err {
		return data, err
	}
	if nil != headers {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
	//请求
	httpResp, err := httpClient.Do(req)
	if err != nil {
		return data, err
	}

	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	data, err = ioutil.ReadAll(httpResp.Body)
	return data, err
}

//IsAjaxRequest 是否是ajax请求
func IsAjaxRequest(r *http.Request) bool {
	return r.Header.Get("X-Requested-With") == "XMLHttpRequest"
}
