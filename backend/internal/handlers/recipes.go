package handlers

import (
	"net/http"

	"backend/internal/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllRecipes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var recipes []models.Recipe

		// Base query with all necessary preloads
		query := db.Preload("Ingredients.Ingredient").
			Preload("Tags.Tag").
			Preload("Steps")

		// Optional filters
		if tag := c.Query("tag"); tag != "" {
			query = query.Joins("JOIN recipe_tags ON recipe_tags.recipe_id = recipes.id").
				Joins("JOIN tags ON tags.id = recipe_tags.tag_id").
				Where("tags.name = ?", tag)
		}

		if ingredient := c.Query("ingredient"); ingredient != "" {
			query = query.Joins("JOIN recipe_ingredients ON recipe_ingredients.recipe_id = recipes.id").
				Joins("JOIN ingredients ON ingredients.id = recipe_ingredients.ingredient_id").
				Where("ingredients.name = ?", ingredient)
		}

		if err := query.Find(&recipes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
			return
		}

		c.JSON(http.StatusOK, recipes)
	}
}

func GetRecipeByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var recipe models.Recipe
		err = db.Preload("Ingredients.Ingredient").Preload("Tags.Tag").Preload("Steps").First(&recipe, id).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipe"})
			return
		}
		c.JSON(http.StatusOK, recipe)
	}
}

func CreateRecipe(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		type StepInput struct {
			Text  string `json:"text"`
			Title string `json:"title"`
		}

		type IngredientInput struct {
			Name   string `json:"name"`
			Amount string `json:"amount"`
			Unit   string `json:"unit"`
		}

		type Input struct {
			Title        string            `json:"title"`
			Instructions string            `json:"instructions"`
			CreatedByID  uint              `json:"created_by_id"`
			Ingredients  []IngredientInput `json:"ingredients"`
			Tags         []string          `json:"tags"`
			Steps        []StepInput       `json:"steps"`
		}

		var input Input
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if len(input.Steps) > 10 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "A recipe can have at most 10 steps."})
			return
		}

		recipe := models.Recipe{
			Title:        input.Title,
			Instructions: input.Instructions,
			CreatedByID:  input.CreatedByID,
		}
		if err := db.Create(&recipe).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create recipe"})
			return
		}

		// Add ingredients
		for _, ing := range input.Ingredients {
			ingredient := models.Ingredient{Name: strings.Title(strings.ToLower(ing.Name))}
			db.FirstOrCreate(&ingredient, ingredient)
			db.Create(&models.RecipeIngredient{
				RecipeID:     recipe.ID,
				IngredientID: ingredient.ID,
				Amount:       ing.Amount,
				Unit:         ing.Unit,
			})
		}

		// Add tags
		for _, tagName := range input.Tags {
			tag := models.Tag{Name: strings.Title(strings.ToLower(tagName))}
			db.FirstOrCreate(&tag, tag)
			db.Create(&models.RecipeTag{
				RecipeID: recipe.ID,
				TagID:    tag.ID,
			})
		}

		for i, step := range input.Steps {
			if strings.TrimSpace(step.Text) == "" {
				continue
			}
			db.Create(&models.RecipeStep{
				RecipeID: recipe.ID,
				Number:   i + 1,
				Title:    strings.TrimSpace(step.Text),
				Text:     strings.TrimSpace(step.Text),
			})
		}

		c.JSON(http.StatusCreated, recipe)
	}
}

func DeleteRecipe(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}
		db.Delete(&models.Recipe{}, id)
		c.Status(http.StatusNoContent)
	}
}

func UpdateRecipe(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		var recipe models.Recipe
		if err := db.First(&recipe, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
			return
		}

		var input struct {
			Title        string `json:"title"`
			Instructions string `json:"instructions"`
			Ingredients  []struct {
				Name   string `json:"name"`
				Amount string `json:"amount"`
				Unit   string `json:"unit"`
			} `json:"ingredients"`
			Tags  []string `json:"tags"`
			Steps []struct {
				Title string `json:"title"`
				Text  string `json:"text"`
			} `json:"steps"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		recipe.Title = input.Title
		recipe.Instructions = input.Instructions
		if err := db.Save(&recipe).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update recipe"})
			return
		}

		// Clear old related data
		db.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeIngredient{})
		db.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeTag{})
		db.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeStep{})

		// Re-create ingredients
		for _, ing := range input.Ingredients {
			var ingredient models.Ingredient
			db.FirstOrCreate(&ingredient, models.Ingredient{Name: strings.Title(strings.ToLower(ing.Name))})

			db.Create(&models.RecipeIngredient{
				RecipeID:     recipe.ID,
				IngredientID: ingredient.ID,
				Amount:       ing.Amount,
				Unit:         ing.Unit,
			})
		}

		// Re-create tags
		for _, tagName := range input.Tags {
			var tag models.Tag
			db.FirstOrCreate(&tag, models.Tag{Name: strings.Title(strings.ToLower(tagName))})

			db.Create(&models.RecipeTag{
				RecipeID: recipe.ID,
				TagID:    tag.ID,
			})
		}

		for i, step := range input.Steps {
			if strings.TrimSpace(step.Text) == "" {
				continue
			}
			db.Create(&models.RecipeStep{
				RecipeID: recipe.ID,
				Number:   i + 1,
				Title:    strings.TrimSpace(step.Title),
				Text:     strings.TrimSpace(step.Text),
			})
		}

		c.JSON(http.StatusOK, gin.H{"message": "Recipe updated"})
	}
}
