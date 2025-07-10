package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"backend/internal/handlers"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestMenuEnv() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.RecipeIngredient{},
		&models.Tag{},
		&models.RecipeTag{},
		&models.MenuEntry{},
		&models.MenuEntryExtra{},
	)
	models.Seed(db)

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Next()
	})
	r.GET("/menu", handlers.GetMenuEntries(db))
	r.POST("/menu", handlers.CreateMenuEntry(db))
	r.DELETE("/menu/:id", handlers.DeleteMenuEntry(db))
	return r, db
}

func TestGetMenuEntries(t *testing.T) {
	r, _ := setupTestMenuEnv()
	req, _ := http.NewRequest("GET", "/menu?user_id=1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
}

func TestCreateMenuEntry(t *testing.T) {
	r, _ := setupTestMenuEnv()

	body := map[string]interface{}{
		"user_id":   1,
		"recipe_id": 1,
		"day":       "Tuesday",
		"meal_type": "lunch",
	}
	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/menu", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.Code)
	}
}

func TestDeleteMenuEntry(t *testing.T) {
	r, db := setupTestMenuEnv()

	entry := models.MenuEntry{UserID: 1, RecipeID: 1, Day: "Friday", MealType: "dinner"}
	db.Create(&entry)

	req, _ := http.NewRequest("DELETE", "/menu/"+strconv.Itoa(int(entry.ID)), nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusNoContent {
		t.Fatalf("expected 204 No Content, got %d", resp.Code)
	}
}
