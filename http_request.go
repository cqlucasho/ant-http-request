package anthttp

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	Glog = log.New(os.Stderr, "[GLog]", log.Ldate)
)

type (
	GLog struct {
		Print func(msg string)
	}

	HttpRequest struct {
		Config		 *Config
		Url          string
		Protocol     string
		Method       string
		Request 	 *http.Request
		RequestData  interface{}
		RequestDataLength int64
		ResponseData []byte
		Retry     	 *Retry

		lock         sync.Mutex
	}
)

/**
 * create a new http request
 */
func NewHttpRequest(requestUrl string, values interface{}) *HttpRequest {
	urlScheme, _ := url.Parse(requestUrl)

	request, _ := http.NewRequest("", "", nil)
	req := &HttpRequest{
		Config: 	 InitialConfig(),
		Url:         requestUrl,
		Protocol:    strings.ToUpper(urlScheme.Scheme),
		Request: 	 request,
		RequestData: values,
		Retry:    	 &Retry{
			retryNum: 0,
			retryMaxNum: 3,
			sleepTime: 1*time.Second,
			Done: make(chan int8),
		},
	}

	return req
}

/**
 * @return *HttpRequest
 */
func (req *HttpRequest) request(body io.ReadCloser) *HttpRequest {
	req.Request.Method = req.Method
	req.Request.URL, req.Request.Host = req.Config.GenerateURL(req.Url)
	req.Config.Client.Transport = req.Config.GetTransportConfig()
	if body != nil {
		req.Request.Body = body
	} else {
		req.Request.Body = req.getBody()
	}

	res, err := req.Config.Client.Do(req.Request)
	if err != nil {
		Glog.Print(err.Error())

		err := TryAgain(req)
		if err != nil {
			Glog.Print(err.Error())
		}

		return req
	}

	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}

	req.ResponseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		Glog.Print(err.Error())
	}

	return req
}

func (req *HttpRequest) getBody() io.ReadCloser {
	if req.RequestData == nil {
		return ioutil.NopCloser(strings.NewReader(""))
	}

	var (
		newReader io.Reader
		dataStr string
	)
	values := reflect.TypeOf(req.RequestData)
	switch values.Kind() {
		case reflect.Map:
			dataStr = req.RequestData.(url.Values).Encode()
			newReader = strings.NewReader(dataStr)
		case reflect.Slice:
			fallthrough
		case reflect.String:
			fallthrough
		default:
			dataStr = req.RequestData.(string)
			newReader = strings.NewReader(dataStr)
	}

	req.RequestDataLength = int64(bytes.Count([]byte(dataStr), nil)-1)
	return ioutil.NopCloser(newReader)
}

func (req *HttpRequest) Get() *HttpRequest {
	req.Method = "GET"

	return req.request(nil)
}

func (req *HttpRequest) Post() *HttpRequest {
	req.Method = "POST"
	req.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req.request(nil)
}

func (req *HttpRequest) PostFile(files []FormFile) error {
	body, err := formFile(files, req)
	if err != nil {
		return err
	}

	req.request(ioutil.NopCloser(body))
	return nil
}

func (req *HttpRequest) Json(method string) *HttpRequest {
	req.Method = strings.ToUpper(method)
	req.Request.Header.Set("Content-Type", "application/json")

	return req.request(nil)
}

func (req *HttpRequest) Xml(method string) *HttpRequest {
	req.Method = strings.ToUpper(method)
	req.Request.Header.Set("Content-Type", "application/xml")

	return req.request(nil)
}

func (req *HttpRequest) SetHeader(key string, value string) *HttpRequest {
	req.Request.Header.Set(key, value)
	return req
}

func (req *HttpRequest) SetHeaders(headers map[string]string) *HttpRequest {
	for k, v := range headers {
		req.Request.Header.Set(k, v)
	}

	return req
}

func (req *HttpRequest) SetCookieJar(jar http.CookieJar) *HttpRequest {
	req.Config.SetCookieJar(jar)

	return req
}

func (req *HttpRequest) SetCheckRedirect(check func(req *http.Request, via []*http.Request) error) *HttpRequest {
	req.Config.SetCheckRedirect(check)

	return req
}

func (req *HttpRequest) SetTLSClientConfig(tlsConfig *tls.Config) *HttpRequest {
	if req.Protocol == "HTTPS" {
		req.Config.SetTLSClientConfig(tlsConfig)
	}

	return req
}

func (req *HttpRequest) SetTransport(transport *http.Transport) *HttpRequest {
	req.Config.SetTransportConfig(transport)

	return req
}

func (req *HttpRequest) SetDialContext(timeout, keepAlive time.Duration) *HttpRequest {
	req.Config.SetDialContext(timeout, keepAlive)

	return req
}

func (req *HttpRequest) SetCerts(certs []tls.Certificate) *HttpRequest {
	req.Config.SetCerts(certs)

	return req
}

func (req *HttpRequest) SetRetry(retry *Retry) *HttpRequest {
	req.Retry = retry

	return req
}

func (req *HttpRequest) SetConfig(config *Config) *HttpRequest {
	req.Config = config

	return req
}
