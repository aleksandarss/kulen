package handlers

import (
	"net/http"
	"strings"

	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetTags(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tags []models.Tag
		if err := db.Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
			return
		}
		c.JSON(http.StatusOK, tags)
	}
}

func CreateTag(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Name string `json:"name"`
		}
		if err := c.ShouldBindJSON(&input); err != nil || strings.TrimSpace(input.Name) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		tag := models.Tag{Name: strings.Title(strings.ToLower(input.Name))}
		if err := db.FirstOrCreate(&tag, tag).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
			return
		}

		c.JSON(http.StatusCreated, tag)
	}
}
