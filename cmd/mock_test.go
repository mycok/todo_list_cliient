package cmd

import (
	"net/http"
	"net/http/httptest"
)

// Mock the todo-list API server responses to use for testing the client.
var testResp = map[string]struct {
	Status int
	Body   string
}{
	"resultsMany": {
		Status: http.StatusOK,
		Body: `{
			"results": [
				{
					"Task": "task 1",
					"Done": false,
					"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
					"CompletedAt": "0001-01-01T00:00:00Z"
				},
				{
					"Task": "task 2",
					"Done": false,
					"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
					"CompletedAt": "0001-01-01T00:00:00Z"
				}
			],
			"date": 356648847899,
			"total_results": 2
		}`,
	},
	"resultsOne": {
		Status: http.StatusOK,
		Body: `{
			"results": [
				{
					"Task": "task 2",
					"Done": false,
					"CreatedAt": "2019-10-28T08:23:38.310097076-04:00",
					"CompletedAt": "0001-01-01T00:00:00Z"
				}
			],
			"date": 356648847899,
			"total_results": 1
		}`,
	},
	"noResults": {
		Status: http.StatusOK,
		Body: `{
			"results": [],
			"date": 356648847899,
			"total_results": 0
		}`,
	},
	"created": {
		Status: http.StatusCreated,
		Body:   "",
	},
	"root": {
		Status: http.StatusOK,
		Body:   "Our API is live",
	},
	"notFound": {
		Status: http.StatusNotFound,
		Body:   "404 - not found",
	},
}

func mockServer(h http.HandlerFunc) (string, func()) {
	s := httptest.NewServer(h)

	return s.URL, func() {
		s.Close()
	}
}
