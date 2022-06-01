package cmd

import (
	"bytes"
	"errors"
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
