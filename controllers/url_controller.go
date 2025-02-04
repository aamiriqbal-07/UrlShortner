package controllers

import (
	"net/http"
	"urlshortner/config"
	"urlshortner/service"

	"github.com/gin-gonic/gin"
)

type URLController struct {
	urlService service.URLService
	config     *config.Config
}

func NewURLController(urlService service.URLService, cfg *config.Config) *URLController {
	return &URLController{
		urlService: urlService,
		config:     cfg,
	}
}

func (c *URLController) ShortenURL(ctx *gin.Context) {
	var request struct {
		URL string `json:"url" binding:"required,url"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL"})
		return
	}

	url, err := c.urlService.ShortenURL(request.URL)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	shortURL := c.config.ShortURL.BaseURL + "/" + url.ShortCode
	ctx.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func (c *URLController) RedirectURL(ctx *gin.Context) {
	shortCode := ctx.Param("shortCode")
	originalURL, err := c.urlService.GetOriginalURL(shortCode)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	ctx.Redirect(http.StatusFound, originalURL)
}

func (c *URLController) GetTopDomains(ctx *gin.Context) {
	metrics, err := c.urlService.GetTopDomains(3)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get metrics"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"domains": metrics})
}
