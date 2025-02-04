package service

import (
    "net/url"
    "strings"
    "urlshortner/config"
    "urlshortner/models"
    "urlshortner/repository"
    "urlshortner/utils"
)

type URLService interface {
    ShortenURL(longURL string) (*models.URL, error)
    GetOriginalURL(shortCode string) (string, error)
    GetTopDomains(limit int) ([]models.DomainMetric, error)
}

type URLServiceImpl struct {
    repo   repository.URLRepository
    config *config.Config
}

func NewURLService(repo repository.URLRepository, cfg *config.Config) URLService {
    return &URLServiceImpl{
        repo:   repo,
        config: cfg,
    }
}

func (s *URLServiceImpl) ShortenURL(longURL string) (*models.URL, error) {
    parsedURL, err := url.Parse(longURL)
    if err != nil {
        return nil, err
    }

    domain := parsedURL.Host
    if strings.HasPrefix(domain, "www.") {
        domain = domain[4:]
    }

    // Check if URL already exists
    if existingURL, err := s.repo.FindByOriginalURL(longURL); err == nil {
        return existingURL, nil
    }

    // Generate new short code
    shortCode := utils.GenerateShortCode(s.config.ShortURL.Length)
    for {
        if _, err := s.repo.FindByShortCode(shortCode); err != nil {
            break
        }
        shortCode = utils.GenerateShortCode(s.config.ShortURL.Length)
    }

    url := &models.URL{
        OriginalURL: longURL,
        ShortCode:   shortCode,
        Domain:      domain,
    }

    if err := s.repo.Create(url); err != nil {
        return nil, err
    }

    return url, nil
}

func (s *URLServiceImpl) GetOriginalURL(shortCode string) (string, error) {
    url, err := s.repo.FindByShortCode(shortCode)
    if err != nil {
        return "", err
    }

    if err := s.repo.IncrementAccessCount(url); err != nil {
        // Log error but don't fail the request
        // logger.Error("Failed to increment access count", err)
    }

    return url.OriginalURL, nil
}

func (s *URLServiceImpl) GetTopDomains(limit int) ([]models.DomainMetric, error) {
    return s.repo.GetTopDomains(limit)
}