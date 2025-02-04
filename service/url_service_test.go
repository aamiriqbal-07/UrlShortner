package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"testing"
	"urlshortner/config"
	"urlshortner/models"
)

type MockURLRepository struct {
	mock.Mock
}

func (m *MockURLRepository) Create(url *models.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepository) FindByShortCode(shortCode string) (*models.URL, error) {
	args := m.Called(shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.URL), args.Error(1)
}

func (m *MockURLRepository) FindByOriginalURL(originalURL string) (*models.URL, error) {
	args := m.Called(originalURL)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.URL), args.Error(1)
}

func (m *MockURLRepository) IncrementAccessCount(url *models.URL) error {
	args := m.Called(url)
	return args.Error(0)
}

func (m *MockURLRepository) GetTopDomains(limit int) ([]models.DomainMetric, error) {
	args := m.Called(limit)
	return args.Get(0).([]models.DomainMetric), args.Error(1)
}

func setupTestService() (*URLServiceImpl, *MockURLRepository) {
	mockRepo := new(MockURLRepository)
	cfg := &config.Config{}
	cfg.ShortURL.Length = 6
	cfg.ShortURL.BaseURL = "http://localhost:8080"
	service := NewURLService(mockRepo, cfg).(*URLServiceImpl)
	return service, mockRepo
}

func TestShortenURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		setupMock   func(*MockURLRepository)
		expectError bool
		expectURL   *models.URL
	}{
		{
			name: "Successfully shorten new URL",
			url:  "https://example.com/page",
			setupMock: func(m *MockURLRepository) {
				m.On("FindByOriginalURL", "https://example.com/page").Return(nil, gorm.ErrRecordNotFound)
				m.On("FindByShortCode", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
				m.On("Create", mock.Anything).Return(nil)
			},
			expectError: false,
		},
		{
			name: "URL already exists",
			url:  "https://example.com/page",
			setupMock: func(m *MockURLRepository) {
				existingURL := &models.URL{
					OriginalURL: "https://example.com/page",
					ShortCode:   "abc123",
					Domain:      "example.com",
				}
				m.On("FindByOriginalURL", "https://example.com/page").Return(existingURL, nil)
			},
			expectError: false,
			expectURL: &models.URL{
				OriginalURL: "https://example.com/page",
				ShortCode:   "abc123",
				Domain:      "example.com",
			},
		},
		{
			name: "Database error on create",
			url:  "https://example.com/page",
			setupMock: func(m *MockURLRepository) {
				m.On("FindByOriginalURL", "https://example.com/page").Return(nil, gorm.ErrRecordNotFound)
				m.On("FindByShortCode", mock.Anything).Return(nil, gorm.ErrRecordNotFound)
				m.On("Create", mock.Anything).Return(errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTestService()
			tt.setupMock(mockRepo)

			url, err := service.ShortenURL(tt.url)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url.ShortCode)
				if tt.expectURL != nil {
					assert.Equal(t, tt.expectURL.ShortCode, url.ShortCode)
					assert.Equal(t, tt.expectURL.OriginalURL, url.OriginalURL)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetOriginalURL(t *testing.T) {
	tests := []struct {
		name        string
		shortCode   string
		setupMock   func(*MockURLRepository)
		expectURL   string
		expectError bool
	}{
		{
			name:      "Successfully get original URL",
			shortCode: "abc123",
			setupMock: func(m *MockURLRepository) {
				url := &models.URL{
					OriginalURL: "https://example.com/page",
					ShortCode:   "abc123",
				}
				m.On("FindByShortCode", "abc123").Return(url, nil)
				m.On("IncrementAccessCount", url).Return(nil)
			},
			expectURL:   "https://example.com/page",
			expectError: false,
		},
		{
			name:      "Short code not found",
			shortCode: "notfound",
			setupMock: func(m *MockURLRepository) {
				m.On("FindByShortCode", "notfound").Return(nil, gorm.ErrRecordNotFound)
			},
			expectError: true,
		},
		{
			name:      "Database error",
			shortCode: "abc123",
			setupMock: func(m *MockURLRepository) {
				m.On("FindByShortCode", "abc123").Return(nil, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTestService()
			tt.setupMock(mockRepo)

			url, err := service.GetOriginalURL(tt.shortCode)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectURL, url)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTopDomains(t *testing.T) {
	tests := []struct {
		name          string
		limit         int
		setupMock     func(*MockURLRepository)
		expectMetrics []models.DomainMetric
		expectError   bool
	}{
		{
			name:  "Successfully get top domains",
			limit: 3,
			setupMock: func(m *MockURLRepository) {
				metrics := []models.DomainMetric{
					{Domain: "example.com", Count: 5},
					{Domain: "test.com", Count: 3},
					{Domain: "demo.com", Count: 1},
				}
				m.On("GetTopDomains", 3).Return(metrics, nil)
			},
			expectMetrics: []models.DomainMetric{
				{Domain: "example.com", Count: 5},
				{Domain: "test.com", Count: 3},
				{Domain: "demo.com", Count: 1},
			},
			expectError: false,
		},
		{
			name:  "Database error",
			limit: 3,
			setupMock: func(m *MockURLRepository) {
				m.On("GetTopDomains", 3).Return([]models.DomainMetric{}, errors.New("database error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, mockRepo := setupTestService()
			tt.setupMock(mockRepo)

			metrics, err := service.GetTopDomains(tt.limit)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectMetrics, metrics)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
