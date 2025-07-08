package handlers

import (
	"net/http"

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
