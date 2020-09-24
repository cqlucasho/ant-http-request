<h1 align="center">Ant</h1>

## Features
  
  * More lightweight, More flexible configuration.
  * Only support GET, POST method, and HTTP/HTTPS.
  * Support upload file.
  * Support retry connection.
  * Support link operations
  
## Usage
  
#### POST

```
testReq := NewHttpRequest("http://127.0.0.1:8000/api/hello", url.Values{"pageIndex": {"1"}, "pageSize":{"1"}})
testReq.Post()
Glog.Print("data: ", string(testReq.ResponseData))
```

#### UPLOAD FILE

```
testReq := NewHttpRequest("http://127.0.0.1:8000/postfile", url.Values{"pageIndex": {"1"}, "pageSize":{"1"}})

var files []FormFile
files = append(files, FormFile{Field: "pf", FilePath: "test.txt"})
files = append(files, FormFile{Field: "pf1", FilePath: "test.txt"})
err := testReq.PostFile(files)
if err != nil {
  Glog.Println(err.Error())
}
```

#### SET CONFIGURATION

```
testReq := NewHttpRequest("http://127.0.0.1:8000/hello", nil)
testReq.SetConfig(&Config{MaxRetryNum: 3}).SetTransport().SetCookieJar().SetRetry().SetTLSClientConfig()
```
