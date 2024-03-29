package httputil

import (
	"net/http"
	"reflect"
	"testing"
	"time"
	"fmt"
	"io/ioutil"
	"net"
	"encoding/json"
)

func TestNewHttpRequest(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want *httpRequest
	}{
		{
			name: "ok",
			args: args{
				url: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewHttpRequest(tt.args.url)
			fmt.Printf("httputil.TestNewhttpRequest() = %v\n", got)
		})
	}
}

func Test_SetHeader(t *testing.T) {
	type fields struct {
		Method      string
		Url         string
		ContentType string
		Header      map[string]string
		Body        map[string]string
		Client      *http.Client
	}
	type args struct {
		header map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpRequest
	}{
		{
			name:   "ok",
			fields: fields{},
			args:   args{header: map[string]string{}},
			want:   &httpRequest{Header: map[string]string{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &httpRequest{
				Url:    tt.fields.Url,
				Header: tt.fields.Header,
				Client: tt.fields.Client,
			}
			if got := h.SetHeader(tt.args.header); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpRequest.SetHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_httpRequest_SetCookies(t *testing.T) {
	type fields struct {
		Url         string
		ContentType string
		Header      map[string]string
		Body        map[string]string
		Cookies     []*http.Cookie
		Client      *http.Client
	}
	type args struct {
		cookies []*http.Cookie
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpRequest
	}{
		{
			name:   "ok",
			fields: fields{},
			args:   args{cookies: []*http.Cookie{}},
			want:   &httpRequest{Cookies: []*http.Cookie{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &httpRequest{
				Url:     tt.fields.Url,
				Header:  tt.fields.Header,
				Cookies: tt.fields.Cookies,
				Client:  tt.fields.Client,
			}
			if got := h.SetCookies(tt.args.cookies); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("httpRequest.SetCookies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetTimeout(t *testing.T) {
	type fields struct {
		Method      string
		Url         string
		ContentType string
		Header      map[string]string
		Body        map[string]string
		Client      *http.Client
	}
	type args struct {
		timeout time.Duration
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *httpRequest
	}{
		{
			name:   "ok",
			fields: fields{Client: &http.Client{}},
			args:   args{timeout: time.Second},
			want: &httpRequest{Client: &http.Client{Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout: time.Second,
				}).DialContext}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &httpRequest{
				Url:    tt.fields.Url,
				Header: tt.fields.Header,
				Client: tt.fields.Client,
			}
			got := h.SetTimeout(tt.args.timeout)
			s, _ := json.Marshal(got)
			fmt.Println("***got:", string(s))
		})
	}
}

func Test_httpRequest_Get(t *testing.T) {
	type fields struct {
		Url     string
		Header  map[string]string
		Body    map[string]string
		Cookies []*http.Cookie
		Client  *http.Client
	}

	tests := []struct {
		name    string
		fields  fields
		want    *http.Response
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Url:     "https://www.baidu.com/",
				Header:  map[string]string{},
				Body:    map[string]string{},
				Cookies: []*http.Cookie{},
				Client: &http.Client{

				},
			},
			want:    &http.Response{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &httpRequest{
				Url:     tt.fields.Url,
				Header:  tt.fields.Header,
				Cookies: tt.fields.Cookies,
				Client:  tt.fields.Client,
			}
			got, err := h.Get(tt.fields.Body)
			bodyData, err := ioutil.ReadAll(got.Body)
			if err != nil {
				panic(err)
			}
			defer got.Body.Close()

			if (err != nil) != tt.wantErr {
				t.Errorf("httpRequest.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println("***resp.body:", string(bodyData))

		})
	}
}
