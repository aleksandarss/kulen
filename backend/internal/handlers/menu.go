package handlers

import (
	"net/http"
	"strconv"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMenuEntries(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
			return
		}

		var entries []models.MenuEntry
		if err := db.
			Where("user_id = ?", userIDStr).
			Preload("Recipe").
			Find(&entries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch menu entries"})
			return
		}

		c.JSON(http.StatusOK, entries)
	}
}

func CreateMenuEntry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			UserID   uint   `json:"user_id"`
			RecipeID uint   `json:"recipe_id"`
			Day      string `json:"day"`
			MealType string `json:"meal_type"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		entry := models.MenuEntry{
			UserID:   input.UserID,
			RecipeID: input.RecipeID,
			Day:      input.Day,
			MealType: input.MealType,
		}
		db.Create(&entry)
		c.JSON(http.StatusCreated, entry)
	}
}

func DeleteMenuEntry(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}
		db.Delete(&models.MenuEntry{}, id)
		c.Status(http.StatusNoContent)
	}
}
