package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
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
			errmsg := fmt.Sprintf("Http server returned status code: %d", resp.StatusCode)
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

func putHttpResource(url string, body string) (string, error) {
	client := NewTimeoutClient(10*time.Second, 15*time.Second)

	reader := strings.NewReader(body)
	req, err := http.NewRequest("PUT", url, reader)

	resp, err := client.Do(req)

	out, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(out), err
}

func deleteHttpResource(url string) (string, error) {
	client := NewTimeoutClient(10*time.Second, 15*time.Second)

	req, err := http.NewRequest("DELETE", url, nil)

	resp, err := client.Do(req)

	out, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(out), err
}

func postHttpResource(url string, body string) (string, error) {
	client := NewTimeoutClient(10*time.Second, 15*time.Second)

	ok := false
	retriesAllowed := 3
	resp := &http.Response{}
	var err error = nil
	bodyType := "application/json"

	reader := strings.NewReader(body)

	for i := 0; i < retriesAllowed; i++ {
		if i > 0 {
			fmt.Println("Attempt", i+1, "of", retriesAllowed)
		}

		resp, err = client.Post(url, bodyType, reader)
		if err != nil {
			fmt.Println("err, postHttpResource, post: ", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			ok = true
			break
		} else {
			errmsg := fmt.Sprintf("Http server returned status code: %s", resp.StatusCode)
			err = errors.New(errmsg)
		}
	}

	out := []byte{}
	if ok {
		out, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("err, getHttpResource, readall: ", err)
		}
	}

	return string(out), err
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
