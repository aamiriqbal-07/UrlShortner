package repository

import (
    "testing"
    "urlshortner/models"
    
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    dsn := "root:12345@tcp(127.0.0.1:3306)/urlshortener_test?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        t.Fatalf("Failed to connect to test database: %v", err)
    }
    
    db.Exec("DROP TABLE IF EXISTS urls")
    db.AutoMigrate(&models.URL{})
    
    return db
}

func TestURLRepository(t *testing.T) {
    db := setupTestDB(t)
    repo := NewURLRepository(db)
    
    t.Run("Create and Find URL", func(t *testing.T) {
        url := &models.URL{
            OriginalURL: "https://example.com",
            ShortCode:   "abc123",
            Domain:      "example.com",
        }
        
        err := repo.Create(url)
        assert.NoError(t, err)
        
        found, err := repo.FindByShortCode("abc123")
        assert.NoError(t, err)
        assert.Equal(t, url.OriginalURL, found.OriginalURL)
    })
}