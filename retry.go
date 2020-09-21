package anthttp

import (
	"context"
	"errors"
	"io/ioutil"
	"time"
)

type Retry struct {
	Done 		chan int8
	timeout 	time.Duration
	sleepTime 	time.Duration
	retryMaxNum	int
	retryNum 	int
}

func TryAgain(req *HttpRequest) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 60 * time.Second)
	req.Request.WithContext(ctx)
	defer cancelFunc()

	go tryConn(req)
	defer close(req.Retry.Done)

	select {
		case done := <-req.Retry.Done:
			if done == 1 {
				return errors.New("done")
			} else {
				return nil
			}
		case <-time.After(60 * time.Second):
			Glog.Print("timeout")
	}

	return errors.New("retry is failed")
}

func tryConn(req *HttpRequest) {
	req.Retry.retryNum++
	if req.Retry.retryNum >= req.Retry.retryMaxNum {
		req.Retry.Done <- 1
	}

	req.Retry.sleepTime <<= 1
	time.Sleep(req.Retry.sleepTime)

	res, err := req.Config.Client.Do(req.Request)
	if err != nil {
		Glog.Println(err.Error())
		tryConn(req)
		return
	}

	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}

	req.ResponseData, err = ioutil.ReadAll(res.Body)
	if err != nil {
		Glog.Print(err.Error())
	}

	req.Retry.Done <- 2
}