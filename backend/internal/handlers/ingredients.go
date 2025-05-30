package handlers

import (
	"backend/internal/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetIngredients(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")

		var ingredients []models.Ingredient
		if err := db.
			Where("LOWER(name) LIKE ?", "%"+strings.ToLower(query)+"%").
			Limit(10).
			Find(&ingredients).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query ingredients"})
			return
		}

		c.JSON(http.StatusOK, ingredients)
	}
}
