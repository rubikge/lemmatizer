package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/rubikge/lemmatizer/internal/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisQueue is a mock implementation of RedisQueue
type MockRedisQueue struct {
	mock.Mock
}

// Ensure MockRedisQueue implements RedisQueueInterface
var _ redis.RedisQueueInterface = (*MockRedisQueue)(nil)

func (m *MockRedisQueue) AddRequestToQueue(data string) (string, error) {
	args := m.Called(data)
	return args.String(0), args.Error(1)
}

func (m *MockRedisQueue) GetResponseFromQueue(taskId string) (string, error) {
	args := m.Called(taskId)
	return args.String(0), args.Error(1)
}

func (m *MockRedisQueue) StartWorker(workerName string, process func(string) (string, error)) error {
	args := m.Called(workerName, process)
	return args.Error(0)
}

func setupTest() (*SearchController, *MockRedisQueue) {
	mockQueue := new(MockRedisQueue)
	controller := NewSearchController(mockQueue)
	return controller, mockQueue
}

func TestProcessHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockResponse   string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful enqueue",
			requestBody:    `{"query": "test"}`,
			mockResponse:   "task123",
			mockError:      nil,
			expectedStatus: 200,
			expectedBody:   map[string]interface{}{"task_id": "task123"},
		},
		{
			name:           "queue error",
			requestBody:    `{"query": "test"}`,
			mockResponse:   "",
			mockError:      fmt.Errorf("queue error"),
			expectedStatus: 500,
			expectedBody:   map[string]interface{}{"error": "Failed to enqueue task"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockQueue := setupTest()
			mockQueue.On("AddRequestToQueue", tt.requestBody).Return(tt.mockResponse, tt.mockError)

			app := fiber.New()
			app.Post("/process", controller.ProcessHandler)

			req := httptest.NewRequest("POST", "/process", bytes.NewBuffer([]byte(tt.requestBody)))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)

			mockQueue.AssertExpectations(t)
		})
	}
}

func TestGetResultHandler(t *testing.T) {
	tests := []struct {
		name           string
		taskID         string
		mockResponse   string
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:   "successful result",
			taskID: "task123",
			mockResponse: `{
				"status": "Success",
				"data": {"result": "test"}
			}`,
			mockError:      nil,
			expectedStatus: 200,
			expectedBody: map[string]interface{}{
				"status": redis.StatusSuccess,
				"data":   map[string]interface{}{"result": "test"},
			},
		},
		{
			name:   "processing result",
			taskID: "task123",
			mockResponse: `{
				"status": "Processing..."
			}`,
			mockError:      nil,
			expectedStatus: 202,
			expectedBody: map[string]interface{}{
				"status": redis.StatusProcessing,
			},
		},
		{
			name:   "error result",
			taskID: "task123",
			mockResponse: `{
				"status": "Error",
				"error": "task failed"
			}`,
			mockError:      nil,
			expectedStatus: 500,
			expectedBody: map[string]interface{}{
				"status": redis.StatusError,
				"error":  "task failed",
			},
		},
		{
			name:           "queue error",
			taskID:         "task123",
			mockResponse:   "",
			mockError:      fmt.Errorf("queue error"),
			expectedStatus: 500,
			expectedBody: map[string]interface{}{
				"status": redis.StatusError,
				"error":  "Internal server error",
			},
		},
		{
			name:   "invalid json response",
			taskID: "task123",
			mockResponse: `{
				invalid json
			}`,
			mockError:      nil,
			expectedStatus: 500,
			expectedBody: map[string]interface{}{
				"status": redis.StatusError,
				"error":  "Invalid response format",
			},
		},
		{
			name:   "missing status field",
			taskID: "task123",
			mockResponse: `{
				"data": "test"
			}`,
			mockError:      nil,
			expectedStatus: 500,
			expectedBody: map[string]interface{}{
				"status": redis.StatusError,
				"error":  "Invalid response format",
			},
		},
		{
			name:   "unknown status",
			taskID: "task123",
			mockResponse: `{
				"status": "Unknown"
			}`,
			mockError:      nil,
			expectedStatus: 500,
			expectedBody: map[string]interface{}{
				"status": redis.StatusError,
				"error":  "Unknown status",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockQueue := setupTest()
			mockQueue.On("GetResponseFromQueue", tt.taskID).Return(tt.mockResponse, tt.mockError)

			app := fiber.New()
			app.Get("/result/:taskID", controller.GetResultHandler)

			req := httptest.NewRequest("GET", "/result/"+tt.taskID, nil)
			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			var responseBody map[string]interface{}
			err = json.Unmarshal(body, &responseBody)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, responseBody)

			mockQueue.AssertExpectations(t)
		})
	}
}
