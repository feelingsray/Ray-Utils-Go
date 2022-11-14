package commonIota

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// httpDo IOTA专用HTTP请求
func httpDo(url, method string, header map[string][]string, body string, timeout int) (string, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s", url), strings.NewReader(body))
	if err != nil {
		return "", err
	}
	if header != nil {
		for k, v := range header {
			req.Header.Add(k, strings.Join(v, ";"))
		}
	}
	req.Header.Set("Content-Type", "application/json")
	//req.SetBasicAuth("r89a0y2p", encode.MD5("0eVr^vjo"))
	req.Header.Add("Connection", "close")
	req.Close = true
	tr := http.Transport{DisableKeepAlives: true}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second, Transport: &tr}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	respBodyByte, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	client.CloseIdleConnections()
	if resp.StatusCode != 200 {
		return string(respBodyByte), errors.New(fmt.Sprintf("[HTTP %d]:%s", resp.StatusCode, string(respBodyByte)))
	} else {
		return string(respBodyByte), nil
	}
}
