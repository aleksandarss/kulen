package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/handlers"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.GET("/recipes", handlers.GetAllRecipes(db))
	r.GET("/recipes/:id", handlers.GetRecipeByID(db))
	r.POST("/recipes", handlers.CreateRecipe(db))
	r.PUT("/recipes/:id", handlers.UpdateRecipe(db))
	r.DELETE("/recipes/:id", handlers.DeleteRecipe(db))
	return r
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.RecipeIngredient{},
		&models.Tag{},
		&models.RecipeTag{},
		&models.MenuEntry{},
		&models.RecipeStep{},
	)
	models.Seed(db)
	return db
}

func TestGetAllRecipes(t *testing.T) {
	db := setupTestDB()
	r := setupTestRouter(db)

	req, _ := http.NewRequest("GET", "/recipes", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
}

func TestSearchRecipesByTitle(t *testing.T) {
	db := setupTestDB()
	r := setupTestRouter(db)

	req, _ := http.NewRequest("GET", "/recipes?search=panc", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var recipes []models.Recipe
	json.NewDecoder(resp.Body).Decode(&recipes)
	assert.Equal(t, 1, len(recipes))
	assert.Equal(t, "Simple Pancakes", recipes[0].Title)
}

func TestSearchRecipesByTag(t *testing.T) {
	db := setupTestDB()
	r := setupTestRouter(db)

	req, _ := http.NewRequest("GET", "/recipes?search=vega", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var recipes []models.Recipe
	json.NewDecoder(resp.Body).Decode(&recipes)
	assert.GreaterOrEqual(t, len(recipes), 1)
}

func TestCreateRecipe(t *testing.T) {
	db := setupTestDB()
	r := setupTestRouter(db)

	payload := `{
		"title": "Tofu Stir Fry",
		"instructions": "Cook tofu with veggies.",
		"created_by_id": 1,
		"ingredients": [
			{ "name": "Tofu", "amount": "200", "unit": "g" },
			{ "name": "Soy Sauce", "amount": "2", "unit": "tbsp" }
		],
		"tags": ["Vegan", "Dinner"]
	}`

	req := httptest.NewRequest("POST", "/recipes", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Check recipe exists
	var recipe models.Recipe
	err := db.Where("title = ?", "Tofu Stir Fry").First(&recipe).Error
	assert.NoError(t, err)
	assert.Equal(t, "Cook tofu with veggies.", recipe.Instructions)

	// Check ingredients were created
	var tofu models.Ingredient
	var soySauce models.Ingredient
	assert.NoError(t, db.Where("name = ?", "tofu").First(&tofu).Error)
	assert.NoError(t, db.Where("name = ?", "soy sauce").First(&soySauce).Error)

	// Check tags were created
	var vegan models.Tag
	var dinner models.Tag
	assert.NoError(t, db.Where("name = ?", "vegan").First(&vegan).Error)
	assert.NoError(t, db.Where("name = ?", "dinner").First(&dinner).Error)

	// Check linking in recipe_ingredients
	var links []models.RecipeIngredient
	err = db.Where("recipe_id = ?", recipe.ID).Find(&links).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(links))

	// Check linking in recipe_tags
	var tagLinks []models.RecipeTag
	err = db.Where("recipe_id = ?", recipe.ID).Find(&tagLinks).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(tagLinks))
}

func TestGetRecipeByID(t *testing.T) {
	db := setupTestDB()
	r := setupTestRouter(db)

	req, _ := http.NewRequest("GET", "/recipes/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
}

func TestUpdateRecipe(t *testing.T) {
	db := setupTestDB()
	r := setupTestRouter(db)

	body := map[string]string{
		"title":        "Updated Title",
		"instructions": "Updated instructions.",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("PUT", "/recipes/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
}

func TestDeleteRecipe(t *testing.T) {
	db := setupTestDB()
	r := setupTestRouter(db)

	req, _ := http.NewRequest("DELETE", "/recipes/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", resp.Code)
	}
}
