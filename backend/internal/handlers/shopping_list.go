package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShoppingItem struct {
	Name   string `json:"name"`
	Unit   string `json:"unit"`
	Amount string `json:"amount"`
}

func GetShoppingList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user_id"})
			return
		}

		var entries []models.MenuEntry
		if err := db.Preload("Recipe.Ingredients.Ingredient").
			Where("user_id = ?", userIDStr).
			Find(&entries).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load menu"})
			return
		}

		type Item struct {
			Name   string
			Unit   string
			Amount float64
		}

		itemMap := make(map[string]*Item)

		for _, entry := range entries {
			for _, ri := range entry.Recipe.Ingredients {
				key := fmt.Sprintf("%s|%s", ri.Ingredient.Name, ri.Unit)

				parsed, err := strconv.ParseFloat(ri.Amount, 64)
				if err != nil {
					continue // skip malformed
				}

				if itemMap[key] == nil {
					itemMap[key] = &Item{
						Name:   ri.Ingredient.Name,
						Unit:   ri.Unit,
						Amount: parsed,
					}
				} else {
					itemMap[key].Amount += parsed
				}
			}
		}

		list := make([]Item, 0, len(itemMap))
		for _, item := range itemMap {
			list = append(list, *item)
		}

		c.JSON(http.StatusOK, list)
	}
}
