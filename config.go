package anthttp

import (
	"crypto/tls"
	"net"
	"net/http"
	urlpkg "net/url"
	"runtime"
	"strings"
	"time"
)

type Config struct {
	Client			*http.Client
	MaxRetryNum  	int8

	transport 		*http.Transport
	tlsClientConfig *tls.Config
}

func InitialConfig() *Config {
	return &Config{
		MaxRetryNum: 5,
		Client: &http.Client{},
		transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   10 * time.Second,
				KeepAlive: 60 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxConnsPerHost: 	   runtime.GOMAXPROCS(2),
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		tlsClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func (c *Config) SetCheckRedirect(check func(req *http.Request, via []*http.Request) error) {
	c.Client.CheckRedirect = check
}

func (c *Config) GetTransportConfig() *http.Transport {
	c.transport.TLSClientConfig = c.tlsClientConfig

	return c.transport
}

func (c *Config) SetTransportConfig(trans *http.Transport) {
	c.transport = trans
}

func (c *Config) GetTlsClientConfig() *tls.Config {
	return c.tlsClientConfig
}

func (c *Config) SetTLSClientConfig(tlsConfig *tls.Config) {
	c.tlsClientConfig = tlsConfig
}

func (c *Config) SetDialContext(timeout, keepAlive time.Duration) {
	dial := (&net.Dialer{
		Timeout: timeout * time.Second,
		KeepAlive: keepAlive * time.Second,
	}).DialContext

	c.transport.DialContext = dial
}

func (c *Config) SetCookieJar(jar http.CookieJar) {
	c.Client.Jar = jar
}

func (c *Config) SetCerts(certs []tls.Certificate) {
	c.tlsClientConfig.Certificates = certs
}

func (c *Config) GenerateRedirectUrl(method, urlStr string, body urlpkg.Values) *http.Request {
	requestUrl, err := http.NewRequest(strings.ToUpper(method), urlStr, strings.NewReader(body.Encode()))
	if err != nil {
		Glog.Println(err)
	}

	return requestUrl
}

func (c *Config) GenerateURL(urlStr string) (url *urlpkg.URL, host string) {
	u, err := urlpkg.Parse(urlStr)
	if err != nil {
		return nil, ""
	}

	if hasPort(host) {
		u.Host = strings.TrimSuffix(host, ":")
	}

	return u, u.Host
}

func hasPort(s string) bool {
	return strings.LastIndex(s, ":") > strings.LastIndex(s, "]")
}

