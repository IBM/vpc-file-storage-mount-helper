/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2023 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMainFunction(t *testing.T) {
	// Positive Test Case
	t.Run("Main Function Success", func(t *testing.T) {
		// Mock the server to return a 200 OK response
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer mockServer.Close()

		os.Setenv("SOCKET_PATH", mockServer.URL)
		go func() {
			defer os.Unsetenv("SOCKET_PATH")
			main()
		}()

		// Test HTTP request to the server
		resp, err := http.Get(mockServer.URL)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestHandleMounting(t *testing.T) {
	sysOp := &MockSystemOperation{}

	// Positive Test Case
	t.Run("Valid Request", func(t *testing.T) {
		jsonBody := `{"stagingTargetPath": "/staging", "targetPath": "/target", "fsType": "ibmshare", "requestID": "123"}`
		req, err := http.NewRequest("POST", "/api/mount", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handleMounting(sysOp)(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Invalid Request
	t.Run("Invalid Request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/mount", strings.NewReader("invalid JSON"))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handleMounting(sysOp)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Mounting Failure
	t.Run("Mounting Failure", func(t *testing.T) {
		sysOp.ExecuteFunc = func(command string, args ...string) (string, error) {
			return "", errors.New("mounting failed")
		}

		jsonBody := `{"stagingTargetPath": "/staging", "targetPath": "/target", "fsType": "ibmshare", "requestID": "123"}`
		req, err := http.NewRequest("POST", "/api/mount", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handleMounting(sysOp)(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandleUnMount(t *testing.T) {
	sysOp := &MockSystemOperation{}

	// Positive Test Case
	t.Run("Valid Request", func(t *testing.T) {
		jsonBody := `{"targetPath": "/target"}`
		req, err := http.NewRequest("POST", "/api/umount", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handleUnMount(sysOp)(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Invalid Request
	t.Run("Invalid Request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/umount", strings.NewReader("invalid JSON"))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handleUnMount(sysOp)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Unmounting Failure
	t.Run("Unmounting Failure", func(t *testing.T) {
		sysOp.ExecuteFunc = func(command string, args ...string) (string, error) {
			return "", errors.New("unmounting failed")
		}

		jsonBody := `{"targetPath": "/target"}`
		req, err := http.NewRequest("POST", "/api/umount", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		handleUnMount(sysOp)(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestDebugLogs(t *testing.T) {
	sysOp := &MockSystemOperation{}

	// Positive Test Case
	t.Run("Valid Request", func(t *testing.T) {
		jsonBody := `{"requestID": "123"}`
		req, err := http.NewRequest("POST", "/api/debugLogs", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		debugLogs(sysOp)(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Invalid Request
	t.Run("Invalid Request", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/api/debugLogs", strings.NewReader("invalid JSON"))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		debugLogs(sysOp)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// journalctl command failed
	t.Run("Log collection failed", func(t *testing.T) {
		sysOp.ExecuteFunc = func(command string, args ...string) (string, error) {
			return "", errors.New("log file creation failed")
		}

		jsonBody := `{"requestID": "123"}`
		req, err := http.NewRequest("POST", "/api/debugLogs", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		debugLogs(sysOp)(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

}

func TestMountStatus(t *testing.T) {
	sysOp := &MockSystemOperation{}

	// Positive Test Case
	t.Run("Valid Request", func(t *testing.T) {
		jsonBody := `{"targetPath": "/target"}`
		req, err := http.NewRequest("GET", "/api/mountStatus", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mountStatus(sysOp)(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Negative Test Case: Invalid Request
	t.Run("Invalid Request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/mountStatus", strings.NewReader("invalid JSON"))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mountStatus(sysOp)(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	// Error Scenario: Find Mount Failure
	t.Run("Find Mount Failure", func(t *testing.T) {
		sysOp.ExecuteFunc = func(command string, args ...string) (string, error) {
			return "", errors.New("find mount failed")
		}

		jsonBody := `{"targetPath": "/target"}`
		req, err := http.NewRequest("GET", "/api/mountStatus", strings.NewReader(jsonBody))
		assert.NoError(t, err)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = req

		mountStatus(sysOp)(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestMountHelperContainerStatus(t *testing.T) {
	t.Run("MountHelperContainerStatus", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		mountHelperContainerStatus(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Mount-helper-container server is live!", response["Message"])
	})
}

// MockSystemOperation is a mock implementation of SystemOperation for testing.
type MockSystemOperation struct {
	ExecuteFunc func(command string, args ...string) (string, error)
}

func (m *MockSystemOperation) Execute(command string, args ...string) (string, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(command, args...)
	}
	return "", nil
}
