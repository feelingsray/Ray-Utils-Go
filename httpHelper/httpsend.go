package httpHelper

import (
	"bytes"
	"io"
	"net/http"

	"github.com/feelingsray/Ray-Utils-Go/encode"
)

func HttpPostWithAuth(url string, body any, username string, password string) (int, string, error) {
	bodyJson, err := encode.DumpJson(body)
	if err != nil {
		return 500, "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(username, password)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(respBody), nil
}

func HttpPost(url string, body any) (int, string, error) {
	bodyJson, err := encode.DumpJson(body)
	if err != nil {
		return 500, "", err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(respBody), nil
}

func HttpGet(url string) (int, string, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, string(respBody), nil
}
