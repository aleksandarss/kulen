package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"backend/internal/handlers"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTagTest() (*gin.Engine, *gorm.DB) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.RecipeIngredient{},
		&models.Tag{},
		&models.RecipeTag{},
		&models.MenuEntry{},
	)
	models.Seed(db)
	r := gin.Default()
	r.GET("/tags", handlers.GetTags(db))
	r.POST("/tags", handlers.CreateTag(db))
	return r, db
}

func TestGetTags(t *testing.T) {
	r, _ := setupTagTest()
	req, _ := http.NewRequest("GET", "/tags", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
}

func TestCreateTag(t *testing.T) {
	r, _ := setupTagTest()
	body := strings.NewReader(`{"name":"Dessert"}`)
	req, _ := http.NewRequest("POST", "/tags", body)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.Code)
	}
}
