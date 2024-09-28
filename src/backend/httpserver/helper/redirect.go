package helper

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/kataras/iris/v12"
)

const (
	E7DataHost = "localhost"
)

func init() {
	redirectClient = NewRedirectClient("localhost:80")
}

var redirectClient *RedirectClient

type RedirectClient struct {
	addr   string
	client *http.Client
}

func NewRedirectClient(addr string) *RedirectClient {
	cli := http.Client{}
	return &RedirectClient{
		addr:   addr,
		client: &cli,
	}
}

func (rc *RedirectClient) Do(req *http.Request) (*http.Response, error) {
	return rc.client.Do(req)
}

func GetHttpRequestFromIris(c iris.Context) *http.Request {
	req := &http.Request{}
	return req
}

func GetHttpResponseForIris(c iris.Context, resp *http.Response) error {
	return nil
}

// HttpRedirect 用于把当前的请求完全请求到另外一个host的某个path下，
func HttpRedirect() func(c iris.Context) {
	return func(c iris.Context) {
		urlProxy, err := url.Parse("http://" + E7DataHost)
		if err != nil {
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(urlProxy)
		proxy.ServeHTTP(c.ResponseWriter(), c.Request())
	}
}
