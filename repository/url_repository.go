package repository

import (
    "gorm.io/gorm"
	"urlshortner/models"
)

type URLRepository interface {
    Create(url *models.URL) error
    FindByShortCode(shortCode string) (*models.URL, error)
    FindByOriginalURL(originalURL string) (*models.URL, error)
    IncrementAccessCount(url *models.URL) error
    GetTopDomains(limit int) ([]models.DomainMetric, error)
}

type URLRepositoryImpl struct {
    db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
    return &URLRepositoryImpl{db: db}
}

func (r *URLRepositoryImpl) Create(url *models.URL) error {
    return r.db.Create(url).Error
}

func (r *URLRepositoryImpl) FindByShortCode(shortCode string) (*models.URL, error) {
    var url models.URL
    err := r.db.Where("short_code = ?", shortCode).First(&url).Error
    return &url, err
}

func (r *URLRepositoryImpl) FindByOriginalURL(originalURL string) (*models.URL, error) {
    var url models.URL
    err := r.db.Where("original_url = ?", originalURL).First(&url).Error
    return &url, err
}

func (r *URLRepositoryImpl) IncrementAccessCount(url *models.URL) error {
    return r.db.Model(url).Update("access_count", gorm.Expr("access_count + ?", 1)).Error
}

func (r *URLRepositoryImpl) GetTopDomains(limit int) ([]models.DomainMetric, error) {
    var metrics []models.DomainMetric
    err := r.db.Model(&models.URL{}).
        Select("domain, COUNT(*) as count").
        Group("domain").
        Order("count DESC").
        Limit(limit).
        Scan(&metrics).Error
    return metrics, err
}