package anthttp

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
)

func TestHttpRequest(t *testing.T) {
	testReq := NewHttpRequest("http://127.0.0.1:8000/api/sku", url.Values{"pageIndex": {"1"}, "pageSize":{"1"}})
	testReq.Post()
	Glog.Print("data: ", string(testReq.ResponseData))
}

func TestSetRedirect(t *testing.T) {
	testReq := NewHttpRequest("http://127.0.0.1:8000/api/sku", url.Values{"pageIndex": {"1"}, "pageSize":{"1"}})
	testReq.SetCheckRedirect(func(req *http.Request, via []*http.Request) error {
		if len(via) >= 1 {
			return errors.New("stop redirect")
		}

		return nil
	})
}

func TestPostFile(t *testing.T) {
	testReq := NewHttpRequest("http://127.0.0.1:8099/postfile", url.Values{"pageIndex": {"1"}, "pageSize":{"1"}})

	var files []FormFile
	files = append(files, FormFile{Field: "pf", FilePath: "test.txt"})
	files = append(files, FormFile{Field: "pf1", FilePath: "test.txt"})
	err := testReq.PostFile(files)
	if err != nil {
		Glog.Println(err.Error())
	}
}
