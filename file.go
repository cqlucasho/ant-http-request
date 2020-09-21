package anthttp

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
)

type FormFile struct {
	Field string
	FilePath string
}

func formFile(files []FormFile, req *HttpRequest) (*bytes.Buffer, error) {
	body, err := addFile(files, req)
	if err != nil {
		Glog.Print(err.Error())
		return nil, err
	}

	return body, nil
}

func addFile(files []FormFile, req *HttpRequest) (*bytes.Buffer, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	defer bodyWriter.Close()

	for _, uploadFile := range files {
		file, err := os.Open(uploadFile.FilePath)
		if err != nil {
			Glog.Panic(err)
		}

		part, err := bodyWriter.CreateFormFile(uploadFile.Field, filepath.Base(uploadFile.FilePath))
		if err != nil {
			Glog.Panic(err)
		}

		_, err = io.Copy(part, file)
		if err != nil {
			Glog.Panic(err)
		}

		file.Close()
	}

	reqData := req.RequestData.(url.Values)
	for k, v := range reqData {
		for _, vv := range v {
			if err := bodyWriter.WriteField(k, vv); err != nil {
				Glog.Panic(err)
			}
		}
	}

	req.SetHeader("Content-Type", bodyWriter.FormDataContentType())
	return bodyBuf, nil
}


