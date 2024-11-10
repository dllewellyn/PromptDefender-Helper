package endpoints

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockImprover struct{}

func (m *MockImprover) Improve(input string) (string, error) {
	if input == "error" {
		return "", errors.New("mock error")
	}
	return "Improved: " + input, nil
}

type MockCache struct {
	data map[string]string
}

func (m *MockCache) Set(ctx context.Context, key, response string) error {
	if key == "error" {
		return errors.New("mock error")
	}
	m.data[key] = response
	return nil
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	if key == "error" {
		return "", errors.New("mock error")
	}
	return m.data[key], nil
}

func (m *MockCache) Exists(ctx context.Context, key string) (bool, error) {
	_, exists := m.data[key]
	return exists, nil
}

func NewMockCache() *MockCache {
	return &MockCache{data: make(map[string]string)}
}

type MockHTMLRender struct{}

func (m *MockHTMLRender) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: nil,
		Name:     name,
		Data:     data,
	}
}

func TestAddImprover(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	ctx := context.Background()
	mockImprover := &MockImprover{}
	mockCache := NewMockCache()
	mockHTMLRender := &MockHTMLRender{}
	router.HTMLRender = mockHTMLRender

	router.LoadHTMLGlob("../templates/*.html") // Pc302

	AddImprover(ctx, router, mockImprover, mockCache)

	t.Run("JSON request", func(t *testing.T) {
		body := map[string]string{"prompt": "This is a test prompt"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/improve", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.Code)
		}

		var response map[string]string
		if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
			t.Fatalf("expected valid JSON response, got error: %v", err)
		}

		if response["response"] != "Improved: This is a test prompt" {
			t.Fatalf("expected 'Improved: This is a test prompt', got '%s'", response["response"])
		}
	})

	t.Run("Form request", func(t *testing.T) {
		formData := "prompt=This is a test prompt"
		req, _ := http.NewRequest("POST", "/improve", bytes.NewBufferString(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.Code)
		}

		if !bytes.Contains(resp.Body.Bytes(), []byte("Improved: This is a test prompt")) {
			t.Fatalf("expected 'Improved: This is a test prompt' in response, got '%s'", resp.Body.String())
		}
	})

	t.Run("Invalid JSON request", func(t *testing.T) {
		invalidJson := "{invalid_json}"
		req, _ := http.NewRequest("POST", "/improve", bytes.NewBufferString(invalidJson))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", resp.Code)
		}
	})

	t.Run("Missing form data", func(t *testing.T) {
		formData := "invalid=data"
		req, _ := http.NewRequest("POST", "/improve", bytes.NewBufferString(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d", resp.Code)
		}
	})

	t.Run("Cache Get error", func(t *testing.T) {
		body := map[string]string{"prompt": "error"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/improve", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", resp.Code)
		}
	})

	t.Run("Improver Improve error", func(t *testing.T) {
		body := map[string]string{"prompt": "error"}
		jsonBody, _ := json.Marshal(body)
		req, _ := http.NewRequest("POST", "/improve", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			t.Fatalf("expected status 500, got %d", resp.Code)
		}
	})

	t.Run("Nil pointer dereference", func(t *testing.T) {
		formData := "prompt=This is a test prompt"
		req, _ := http.NewRequest("POST", "/improve", bytes.NewBufferString(formData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", resp.Code)
		}

		if !bytes.Contains(resp.Body.Bytes(), []byte("Improved: This is a test prompt")) {
			t.Fatalf("expected 'Improved: This is a test prompt' in response, got '%s'", resp.Body.String())
		}
	})
}
