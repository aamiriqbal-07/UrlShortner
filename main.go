package main

import (
	"log"
	"urlshortner/config"
	"urlshortner/controllers"
	"urlshortner/models"
	"urlshortner/repository"
	"urlshortner/service"

	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupRouter(controller *controllers.URLController) *gin.Engine {
	router := gin.Default()

	router.POST("/api/v1/shorten", controller.ShortenURL)
	router.GET("/:shortCode", controller.RedirectURL)
	router.GET("/api/v1/metrics/top-domains", controller.GetTopDomains)

	return router
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Setup database connection
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.URL{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, cfg)
	urlController := controllers.NewURLController(urlService, cfg)

	router := setupRouter(urlController)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
