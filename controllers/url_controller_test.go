package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortner/config"
	"urlshortner/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) ShortenURL(longURL string) (*models.URL, error) {
	args := m.Called(longURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.URL), args.Error(1)
}

func (m *MockURLService) GetOriginalURL(shortCode string) (string, error) {
	args := m.Called(shortCode)
	return args.String(0), args.Error(1)
}

func (m *MockURLService) GetTopDomains(limit int) ([]models.DomainMetric, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.DomainMetric), args.Error(1)
}

func setupTestController() (*URLController, *MockURLService, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockURLService)
	cfg := &config.Config{}
	cfg.ShortURL.BaseURL = "http://localhost:8080"
	controller := NewURLController(mockService, cfg)
	router := gin.New()
	return controller, mockService, router
}

func TestShortenURLEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMock      func(*MockURLService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successfully shorten URL",
			requestBody: map[string]interface{}{
				"url": "https://example.com/page",
			},
			setupMock: func(m *MockURLService) {
				m.On("ShortenURL", "https://example.com/page").Return(&models.URL{
					OriginalURL: "https://example.com/page",
					ShortCode:   "abc123",
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"short_url": "http://localhost:8080/abc123",
			},
		},
		{
			name: "Invalid URL format",
			requestBody: map[string]interface{}{
				"url": "not-a-valid-url",
			},
			setupMock: func(m *MockURLService) {
				// No mock needed - validation will fail
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid URL",
			},
		},
		{
			name: "Service error",
			requestBody: map[string]interface{}{
				"url": "https://example.com/page",
			},
			setupMock: func(m *MockURLService) {
				m.On("ShortenURL", "https://example.com/page").Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to shorten URL",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, router := setupTestController()
			tt.setupMock(mockService)

			router.POST("/api/v1/shorten", controller.ShortenURL)

			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest("POST", "/api/v1/shorten", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestRedirectEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		shortCode      string
		setupMock      func(*MockURLService)
		expectedStatus int
		expectedURL    string
	}{
		{
			name:      "Successful redirect",
			shortCode: "abc123",
			setupMock: func(m *MockURLService) {
				m.On("GetOriginalURL", "abc123").Return("https://example.com/page", nil)
			},
			expectedStatus: http.StatusFound,
			expectedURL:    "https://example.com/page",
		},
		{
			name:      "Short code not found",
			shortCode: "notfound",
			setupMock: func(m *MockURLService) {
				m.On("GetOriginalURL", "notfound").Return("", errors.New("not found"))
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, router := setupTestController()
			tt.setupMock(mockService)

			router.GET("/:shortCode", controller.RedirectURL)

			req := httptest.NewRequest("GET", "/"+tt.shortCode, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedURL != "" {
				assert.Equal(t, tt.expectedURL, w.Header().Get("Location"))
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetTopDomainsEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockURLService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "Successfully get top domains",
			setupMock: func(m *MockURLService) {
				metrics := []models.DomainMetric{
					{Domain: "example.com", Count: 5},
					{Domain: "test.com", Count: 3},
					{Domain: "demo.com", Count: 1},
				}
				m.On("GetTopDomains", 3).Return(metrics, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"domains": []interface{}{
					map[string]interface{}{"domain": "example.com", "count": float64(5)},
					map[string]interface{}{"domain": "test.com", "count": float64(3)},
					map[string]interface{}{"domain": "demo.com", "count": float64(1)},
				},
			},
		},
		{
			name: "Service returns empty result",
			setupMock: func(m *MockURLService) {
				m.On("GetTopDomains", 3).Return([]models.DomainMetric{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"domains": []interface{}{},
			},
		},
		{
			name: "Service error",
			setupMock: func(m *MockURLService) {
				m.On("GetTopDomains", 3).Return([]models.DomainMetric{}, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to get metrics",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controller, mockService, router := setupTestController()
			tt.setupMock(mockService)

			router.GET("/api/v1/metrics/top-domains", controller.GetTopDomains)

			req := httptest.NewRequest("GET", "/api/v1/metrics/top-domains", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, tt.expectedBody, response)

			mockService.AssertExpectations(t)
		})
	}
}

func TestInvalidJSONRequest(t *testing.T) {
	controller, _, router := setupTestController()
	router.POST("/api/v1/shorten", controller.ShortenURL)

	// Send invalid JSON
	req := httptest.NewRequest("POST", "/api/v1/shorten", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMissingURLParameter(t *testing.T) {
	controller, _, router := setupTestController()
	router.POST("/api/v1/shorten", controller.ShortenURL)

	// Send JSON without URL parameter
	body, _ := json.Marshal(map[string]interface{}{
		"not_url": "https://example.com",
	})
	req := httptest.NewRequest("POST", "/api/v1/shorten", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
