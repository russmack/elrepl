package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func getHttpResource(url string) (string, error) {
	client := NewTimeoutClient(10*time.Second, 15*time.Second)

	ok := false
	retriesAllowed := 3
	resp := &http.Response{}
	var err error = nil

	for i := 0; i < retriesAllowed; i++ {
		if i > 0 {
			fmt.Println("Attempt", i+1, "of", retriesAllowed)
		}

		resp, err = client.Get(url)
		if err != nil {
			fmt.Println("err, getHttpResource, get: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			ok = true
			break
		} else {
			errmsg := fmt.Sprintf("Http server returned status code:", resp.StatusCode)
			err = errors.New(errmsg)
		}
	}

	body := []byte{}
	if ok {
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("err, getHttpResource, readall: ", err)
		}
	}

	return string(body), err
}

func TimeoutDialer(connectTimeout time.Duration, readWriteTimeout time.Duration) func(network, address string) (conn net.Conn, err error) {
	return func(network, address string) (net.Conn, error) {
		conn, err := net.DialTimeout(network, address, connectTimeout)
		if err != nil {
			return nil, err
		}
		err = conn.SetDeadline(time.Now().Add(readWriteTimeout))
		if err != nil {
			fmt.Println("Unable to set deadline of http connection:", err)
			return nil, err
		}
		return conn, nil
	}
}

func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}
}
