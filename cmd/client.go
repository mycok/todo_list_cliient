package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrConnection = errors.New("connection error")
	ErrNotFound = errors.New("not found")
	ErrInvalidResponse = errors.New("invalid response")
	ErrInvalid = errors.New("invalid data")
	ErrNotNumber = errors.New("not a number")
)

type item struct {
	Task string
	Done bool
	CreatedAt time.Time
	CompletedAt time.Time
}

type response struct {
	Results []item `json:"results"`
	Date int `json:"date"`
	TotalResults int `json:"total_results"`
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return c
}

func getItems(url string) ([]item, error) {
	resp, err := newClient().Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}

		err = ErrInvalidResponse

		if resp.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}

		return nil, fmt.Errorf("%w: %s", err, errMsg)
	}

	var respData response

	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	if len(respData.Results) == 0 {
		return nil, fmt.Errorf("%w: no results found", ErrNotFound)
	}

	return respData.Results, nil
}

func getAll(url string) ([]item, error) {
	u := fmt.Sprintf("%s/todo", url)

	return getItems(u)
}
