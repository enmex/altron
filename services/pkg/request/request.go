package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func Post[T any](url string, request interface{}) (*T, error) {
	requestJson, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(requestJson),
	)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", body.Message)
	}

	var body T
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func PostWithEmptyResponse(url string, request interface{}) (int, error) {
	requestJson, err := json.Marshal(request)

	if err != nil {
		return 0, err
	}

	res, err := http.Post(
		url,
		"application/json",
		bytes.NewBuffer(requestJson),
	)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("%s", body.Message)
	}

	return res.StatusCode, nil
}

func Get[T any](url string) (*T, error) {
	res, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", body.Message)
	}

	var body T
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func DefaultGet(url string) ([]byte, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, 0, err
		}
		return nil, 0, fmt.Errorf("%s", body.Message)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}

	return body, res.StatusCode, nil
}

func Put(url string, request interface{}) error {
	requestJson, err := json.Marshal(request)
	if err != nil {
		return err
	}

	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestJson))
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return err
		}
		return fmt.Errorf("%s", body.Message)
	}

	return nil
}

func EmptyPut(url string) error {
	client := &http.Client{}

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return err
		}
		return fmt.Errorf("%s", body.Message)
	}

	return nil
}

func Delete[T any](url string) (*T, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", body.Message)
	}

	var body T
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func DeleteWithEmptyResponse(url string) (int, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, nil)

	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return 0, err
		}
		return res.StatusCode, fmt.Errorf("%s", body.Message)
	}

	return res.StatusCode, nil
}

func Patch[T any](url string, request interface{}) (*T, error) {
	requestJson, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	client := http.Client{}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(requestJson))

	if err != nil {
		return nil, nil
	}

	res, err := client.Do(req)

	if err != nil {
		return nil, nil
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", body.Message)
	}

	var body T
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return &body, nil
}

func PatchWithEmptyResponse(url string, request interface{}) (int, error) {
	requestJson, err := json.Marshal(request)

	if err != nil {
		return 0, err
	}

	client := http.Client{}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(requestJson))

	if err != nil {
		return 0, nil
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("%s", body.Message)
	}

	return res.StatusCode, nil
}

func DeleteWithRequest(url string, request interface{}) (int, error) {
	requestJson, err := json.Marshal(request)

	if err != nil {
		return 0, err
	}

	client := http.Client{}

	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(requestJson))

	if err != nil {
		return 0, err
	}

	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		var body ErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			return 0, err
		}
		return 0, fmt.Errorf("%s", body.Message)
	}

	return res.StatusCode, nil
}
