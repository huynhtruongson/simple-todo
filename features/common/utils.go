package common

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func MakeHttpRequest(url, method string, payload []byte, acToken string) (int, []byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", acToken))

	if err != nil {
		return 0, nil, fmt.Errorf("make http request error,%w\n", err)
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return 0, nil, fmt.Errorf("call create user api error,%w\n", err)
	}
	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("could not read response body: %w\n", err)
	}
	return resp.StatusCode, resBody, nil
}