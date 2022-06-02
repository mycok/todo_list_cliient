package cmd

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

func TestListAction(t *testing.T) {
	testCases := []struct {
		name           string
		expectedErr    error
		expectedOutput string
		resp           struct {
			Status int
			Body   string
		}
		closeServer bool
	}{
		{
			name:           "WithResults",
			expectedErr:    nil,
			expectedOutput: "𝘅  1  task 1\n𝘅  2  task 2\n",
			resp:           testResp["resultsMany"],
		},
		{
			name:        "NoResults",
			expectedErr: ErrNotFound,
			resp:        testResp["noResults"],
		},
		{
			name:        "InvalidURL",
			expectedErr: ErrConnection,
			resp:        testResp["noResults"],
			closeServer: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				w.Write([]byte(tc.resp.Body))
				// fmt.Fprintln(w, tc.resp.Body)
			})

			defer cleanup()

			if tc.closeServer {
				cleanup()
			}

			var outputBuf bytes.Buffer

			err := listAction(&outputBuf, url)

			// Handle the error path.
			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Expected error: %q, but got no error instead", tc.expectedErr)

				}

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected error: %q, but got: %q instead", tc.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected no error, but got: %q instead", err)
			}

			// Handle success path.
			if tc.expectedOutput != outputBuf.String() {
				t.Errorf("Expected output: %s, but got: %s", tc.expectedOutput, outputBuf.String())
			}
		})
	}
}

func TestViewAction(t *testing.T) {
	testCases := []struct {
		name           string
		expectedErr    error
		expectedOutput string
		resp           struct {
			Status int
			Body   string
		}
		id string
	}{
		{
			name:        "WithSingleResult",
			expectedErr: nil,
			expectedOutput: "Task:         task 2\n" +
				"Created at:   Oct/28 @08:00\n" +
				"Completed:    No\n",
			resp: testResp["resultsOne"],
			id:   "2",
		},
		{
			name:        "NoFound",
			expectedErr: ErrNotFound,
			resp:        testResp["noResults"],
			id:          "1",
		},
		{
			name:        "InvalidID",
			expectedErr: ErrNotNumber,
			resp:        testResp["noResults"],
			id:          "me",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.resp.Status)
				w.Write([]byte(tc.resp.Body))
			})

			defer cleanup()

			var outputBuf bytes.Buffer

			err := viewAction(&outputBuf, url, tc.id)

			// Handle the error path.
			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("Expected error: %q, but got no error instead", tc.expectedErr)

				}

				if !errors.Is(err, tc.expectedErr) {
					t.Errorf("Expected error: %q, but got: %q instead", tc.expectedErr, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Expected no error, but got: %q instead", err)
			}

			// Handle success path.
			if tc.expectedOutput != outputBuf.String() {
				t.Errorf("Expected output: %s\n, but got: %s", tc.expectedOutput, outputBuf.String())
			}
		})
	}
}

func TestAddAction(t *testing.T) {
	expectedURLPath := "/todo"
	expectedMethod := http.MethodPost
	expectedBody := "{\"task\":\"task 1\"}\n"
	expectedContentType := "application/json"
	args := []string{"task", "1"}
	expectedOutput := "Added item: task 1 : to the list\n"

	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedURLPath {
			t.Fatalf("Expected path: %s, but got: %s instead", expectedURLPath, r.URL.Path)
		}

		if r.Method != expectedMethod {
			t.Fatalf("Expected http method: %s, but got: %s instead", expectedMethod, r.Method)
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		r.Body.Close()

		if string(body) != expectedBody {
			t.Errorf("Expected body: %s, but got: %s instead", expectedBody, string(body))
		}

		if r.Header.Get("Content-Type") != expectedContentType {
			t.Fatalf("Expected content-type: %s, but got: %s instead", expectedContentType, r.Header.Get("Content-Type"))
		}

		w.Header().Set("Content-Type", expectedContentType)
		w.WriteHeader(testResp["created"].Status)
		w.Write([]byte(testResp["created"].Body))
	})

	defer cleanup()

	var body bytes.Buffer

	if err := addAction(&body, url, args); err != nil {
		t.Fatalf("Expected no error, but got: %q instead", err)
	}

	if expectedOutput != body.String() {
		t.Errorf("Expected output: %s, but got: %s instead", expectedOutput, body.String())
	}
}

func TestCompleteAction(t *testing.T) {
	expectedURLPath := "/todo/1"
	expectedMethod := http.MethodPatch
	expectedQuery := "complete"
	expectedOutput := "Item number 1 marked as complete\n"
	arg := "1"

	url, cleanup := mockServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != expectedURLPath {
			t.Fatalf("Expected path: %s, but got: %s instead", expectedURLPath, r.URL.Path)
		}

		if r.Method != expectedMethod {
			t.Fatalf("Expected http method: %s, but got: %s instead", expectedMethod, r.Method)
		}

		if _, ok := r.URL.Query()[expectedQuery]; !ok {
			t.Fatalf("Expected path query: %s not found in URL", expectedQuery)
		}

		w.WriteHeader(testResp["noContent"].Status)
		w.Write([]byte(testResp["noContent"].Body))
	})

	defer cleanup()

	var body bytes.Buffer

	if err := completeAction(&body, url, arg); err != nil {
		t.Fatalf("Expected no error, but got: %q instead", err)
	}

	if expectedOutput != body.String() {
		t.Errorf("Expected output: %s, but got: %s instead", expectedOutput, body.String())
	}
}

