package requests

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func request(method string, url string, data map[string]string) (string, error) {
	dataRaw := ""
	for k, v := range data {
		dataRaw += k + "=" + v + "&"
	}

	var (
		req *http.Request
		err error
	)
	if method == http.MethodPost {
		req, err = http.NewRequest(method, url, strings.NewReader(dataRaw))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req, err = http.NewRequest(method, url+"?"+dataRaw, nil)
	}
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Post(url string, data map[string]string) (string, error) {
	return request(http.MethodPost, url, data)
}

func Get(url string, data map[string]string) (string, error) {
	return request(http.MethodGet, url, data)
}
