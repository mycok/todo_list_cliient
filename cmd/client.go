package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const timeFormat = "Jan/02 @15:00"

var (
	ErrConnection      = errors.New("connection error")
	ErrNotFound        = errors.New("not found")
	ErrInvalidResponse = errors.New("invalid response")
	ErrInvalid         = errors.New("invalid data")
	ErrNotNumber       = errors.New("not a number")
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type response struct {
	Results      []item `json:"results"`
	Date         int    `json:"date"`
	TotalResults int    `json:"total_results"`
}

func getAll(url string) ([]item, error) {
	u := fmt.Sprintf("%s/todo", url)

	return getItems(u)
}

func getItem(url string, id int) (item, error) {
	u := fmt.Sprintf("%s/todo/%d", url, id)

	items, err := getItems(u)
	if err != nil {
		return item{}, err
	}

	return items[0], nil
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

func addItem(url, name string) error {
	u := fmt.Sprintf("%s/todo", url)

	var body bytes.Buffer

	item := struct {
		Task string `json:"task"`
	}{
		Task: name,
	}

	if err := json.NewEncoder(&body).Encode(item); err != nil {
		return err
	}

	return sendMutatingRequest(u, http.MethodPost, "application/json", http.StatusCreated, &body)
}

func completeItem(url string, id int) error {
	u := fmt.Sprintf("%s/todo/%d?complete", url, id)

	return sendMutatingRequest(u, http.MethodPatch, "", http.StatusNoContent, nil)
}

func deleteItem(url string, id int) error {
	u := fmt.Sprintf("%s/todo/%d", url, id)

	return sendMutatingRequest(u, http.MethodDelete, "", http.StatusNoContent, nil)
}

func sendMutatingRequest(url, method, contentType string, statusCode int, body io.Reader) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := newClient().Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != statusCode {
		msg, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %w", err)
		}

		err = ErrInvalidResponse

		if resp.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}

		return fmt.Errorf("%w: %s", err, msg)
	}

	return nil
}

func newClient() *http.Client {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	return c
}
