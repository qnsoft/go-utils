package httputil

import (
	"net/http"
	"net"
	"time"
	"strings"
	"net/url"
	"encoding/json"
	"os"
	"mime/multipart"
	"bytes"
	"io"
	"fmt"
)

const (
	FormContentType      = "application/x-www-form-urlencoded"
	JsonContentType      = "application/json"
	MultipartContentType = "multipart/form-data"
)

// http请求类
type httpRequest struct {
	// 网关
	Url string `json:"url"`
	// 请求头
	Header map[string]string `json:"header"`
	// cookie数组
	Cookies []*http.Cookie `json:"cookies"`
	// http客户端
	Client *http.Client `json:"client"`
}

// 创建http请求
func NewHttpRequest(url string) *httpRequest {
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: 2 * time.Second,
		}).DialContext,
	}

	client := &http.Client{Transport: transport}
	return &httpRequest{
		Url:    url,
		Client: client,
	}
}

// 设置请求头
func (h *httpRequest) SetHeader(header map[string]string) *httpRequest {
	h.Header = header
	return h
}

// 设置请求超时时间，这里的超时指的是连接超时
func (h *httpRequest) SetTimeout(timeout time.Duration) *httpRequest {
	h.Client.Transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: timeout,
		}).DialContext,
	}
	return h
}

// 设置cookies
func (h *httpRequest) SetCookies(cookies []*http.Cookie) *httpRequest {
	h.Cookies = cookies
	return h
}

// Get请求
func (h *httpRequest) Get(data map[string]string) (*http.Response, error) {
	return h.request(http.MethodGet, FormContentType, handleFormData(data))
}

// Post json请求
func (h *httpRequest) Post(data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return h.request(http.MethodPost, JsonContentType, strings.NewReader(string(jsonData)))
}

// Post表单请求
func (h *httpRequest) PostForm(data map[string]string) (*http.Response, error) {
	return h.request(http.MethodPost, FormContentType, handleFormData(data))
}

// 处理表单数据
func handleFormData(data map[string]string) io.Reader {
	values := url.Values{}
	for k, v := range data {
		values.Add(k, v)
	}
	return strings.NewReader(values.Encode())
}

// 文件上传
func (h *httpRequest) UploadFile(fileName, path string, data map[string]string) (*http.Response, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	formFileW, err := w.CreateFormFile("fileName", fileName)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(formFileW, file)
	if err != nil {
		return nil, err
	}

	for k, v := range data {
		w.WriteField(k, fmt.Sprint(v))
	}
	return h.request(http.MethodPost, MultipartContentType, body)
}

// 请求
func (h *httpRequest) request(method string, contentType string, body io.Reader) (*http.Response, error) {
	// 创建请求
	req, err := http.NewRequest(method, h.Url, body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	// 添加请求头
	for k, v := range h.Header {
		if strings.ToLower(k) == "host" {
			req.Host = fmt.Sprint(v)
		} else {
			req.Header.Add(k, fmt.Sprint(v))
		}
	}
	req.Header.Add("Content-Type", contentType)
	// 添加Cookie
	for _, v := range h.Cookies {
		req.AddCookie(v)
	}
	return h.Client.Do(req)
}
