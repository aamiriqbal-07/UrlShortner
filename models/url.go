package models

import (
    "time"
)

type URL struct {
    ID          uint      `gorm:"primarykey"`
    OriginalURL string    `gorm:"type:text;not null"`
    ShortCode   string    `gorm:"type:varchar(10);uniqueIndex;not null"`
    Domain      string    `gorm:"type:varchar(255);index;not null"`
    CreatedAt   time.Time
    AccessCount int       `gorm:"default:0"`
}

type DomainMetric struct {
    Domain string `json:"domain"`
    Count  int    `json:"count"`
}