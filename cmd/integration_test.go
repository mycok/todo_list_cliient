//go:build integration
// +build integration

package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {
	apiRoot := "http://localhost:8080"

	if os.Getenv("TODO_API_ROOT") != "" {
		apiRoot = os.Getenv("TODO_API_ROOT")
	}

	today := time.Now().Format("Jan/02")

	tName := randomTaskName(t)

	taskID := ""

	// Integration tests to execute include
	// [Add, List, View, Complete, ListComplete, Delete, ListDeletedTask].
	t.Run("Add", func(t *testing.T) {
		outputBuf := &bytes.Buffer{} 

		args := []string{tName}

		expectedOutput := fmt.Sprintf("Added item: %s : to the list\n", tName)

		if err := addAction(outputBuf, apiRoot, args); err != nil {
			t.Fatalf("Expected no error, but got: %q instead", err)
		}

		if expectedOutput != outputBuf.String() {
			t.Errorf(
				"Expected output: %s, but got: %s instead",
				expectedOutput,
				outputBuf.String(),
			)
		}
	})

	t.Run("List", func(t *testing.T) {
		outputBuf := &bytes.Buffer{}

		if err := listAction(outputBuf, apiRoot); err != nil {
			t.Fatalf("Expected no error, but got: %q instead", err)
		}

		// Single line of text representing a todo type from the
		//  listAction output.
		outputTask := ""

		scanner := bufio.NewScanner(outputBuf)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), tName) {
				outputTask = scanner.Text()

				break
			}
		}

		if outputTask == "" {
			t.Errorf("Task %s is not in the list", tName)
		}

		taskCompletionStatus := strings.Fields(outputTask)[0]
		if taskCompletionStatus != "𝘅" {
			t.Errorf(
				"Expected status to be: %s, but got: %s instead",
				"𝘅",
				taskCompletionStatus,
			)
		}

		// Set the taskID variable using the newly added task ID.
		taskID = strings.Fields(outputTask)[1]
	})

	vResp := t.Run("View", func(t *testing.T) {
		outputBuf := &bytes.Buffer{}

		if err := viewAction(outputBuf, apiRoot, taskID); err != nil {
			t.Fatalf("Expected no error, but got: %q instead", err)
		}

		expectedOutput := strings.Split(outputBuf.String(), "\n")

		if !strings.Contains(expectedOutput[0], tName) {
			t.Errorf(
				"Expected task: %s, but got: %s instead",
				tName,
				expectedOutput[0],
			)
		}

		if !strings.Contains(expectedOutput[1], today) {
			t.Errorf(
				"Expected creation date: %s, but got: %s instead",
				today,
				expectedOutput[1],
			)
		}

		if !strings.Contains(expectedOutput[2], "No") {
			t.Errorf(
				"Expected status: %s, but got: %s instead",
				"No",
				expectedOutput[2],
			)
		}
	})

	// If the view test fails, stop all test execution at this point
	// by calling t.Fatal() from the main test function.
	// Successive test cases require a successful view test response
	//  to continue.
	if !vResp {
		t.Fatal("View task test failed: terminating integration tests.")
	}

	t.Run("Complete", func(t *testing.T) {
		outputBuf := &bytes.Buffer{}

		if err := completeAction(outputBuf, apiRoot, taskID); err != nil {
			t.Fatalf("Expected no error, but got: %q instead", err)
		}

		expectedOutput := fmt.Sprintf(
			"Item number %s marked as complete\n",
			taskID,
		)

		if expectedOutput != outputBuf.String() {
			t.Errorf(
				"Expected output: %s, but got: %s instead",
				expectedOutput,
				outputBuf.String(),
			)
		}
	})

	t.Run("ListAfterComplete", func(t *testing.T) {
		outputBuf := &bytes.Buffer{}

		if err := listAction(outputBuf, apiRoot); err != nil {
			t.Fatalf("Expected no error, but got: %q instead", err)
		}

		// Single line of text representing a todo type from the
		//  listAction output.
		outputTask := ""

		scanner := bufio.NewScanner(outputBuf)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), tName) {
				outputTask = scanner.Text()

				break
			}
		}

		if outputTask == "" {
			t.Errorf("Task %s is not in the list", tName)
		}

		taskCompletionStatus := strings.Fields(outputTask)[0]
		if taskCompletionStatus != "✅" {
			t.Errorf(
				"Expected status to be: %s, but got: %s instead",
				"✅",
				taskCompletionStatus,
			)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		outputBuf := &bytes.Buffer{}

		if err := deleteAction(outputBuf, apiRoot, taskID); err != nil {
			t.Fatalf("Expected no error, but got: %q instead", err)
		}

		expectedOutput := fmt.Sprintf("Item number %s deleted from the list\n", taskID)

		if expectedOutput != outputBuf.String() {
			t.Errorf(
				"Expected output: %s, but got: %s instead",
				expectedOutput,
				outputBuf.String(),
			)
		}
	})

	t.Run("ListAfterDelete", func(t *testing.T) {
		outputBuf := &bytes.Buffer{}

		if err := listAction(outputBuf, apiRoot); err != nil {
			t.Fatalf("Expected no error, but got: %q instead", err)
		}

		scanner := bufio.NewScanner(outputBuf)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), tName) {
				t.Errorf("Task: %s is still in the list", tName)

				break
			}
		}
	})
}

func randomTaskName(t *testing.T) string {
	t.Helper()

	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMOPQRSTUVWXYZ12345678910"

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var p strings.Builder

	for i := 0; i < 32; i++ {
		p.WriteByte(chars[r.Intn(len(chars))])
	}

	return p.String()
}
