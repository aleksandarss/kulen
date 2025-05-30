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
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDVal.(uint)

		var entries []models.MenuEntry
		if err := db.
			Where("user_id = ?", userID).
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
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDVal.(uint)

		var input struct {
			RecipeID uint   `json:"recipe_id"`
			Day      string `json:"day"`
			MealType string `json:"meal_type"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// Remove existing entry for the same user/day/meal
		db.Where("user_id = ? AND day = ? AND meal_type = ?", userID, input.Day, input.MealType).
			Delete(&models.MenuEntry{})

		entry := models.MenuEntry{
			UserID:   userID,
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
		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		userID := userIDVal.(uint)

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
			return
		}

		var entry models.MenuEntry
		if err := db.Where("id = ? AND user_id = ?", id, userID).First(&entry).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
			return
		}

		db.Delete(&entry)
		c.Status(http.StatusNoContent)
	}
}
