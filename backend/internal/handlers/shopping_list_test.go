package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/handlers"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestShoppingEnv() (*gin.Engine, *gorm.DB) {
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
	r.GET("/shopping-list", handlers.GetShoppingList(db))
	return r, db
}

func TestGetShoppingList(t *testing.T) {
	r, _ := setupTestShoppingEnv()

	req, _ := http.NewRequest("GET", "/shopping-list?user_id=1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
}

func TestGetShoppingListInvalidUser(t *testing.T) {
	r, _ := setupTestShoppingEnv()

	req, _ := http.NewRequest("GET", "/shopping-list?user_id=abc", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
}

func TestGetShoppingListMissingUser(t *testing.T) {
	r, _ := setupTestShoppingEnv()

	req, _ := http.NewRequest("GET", "/shopping-list", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request, got %d", resp.Code)
	}
}

func TestShoppingListEmpty(t *testing.T) {
	r, db := setupTestShoppingEnv()

	db.Where("1 = 1").Delete(&models.MenuEntry{})

	req, _ := http.NewRequest("GET", "/shopping-list?user_id=1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	if resp.Body.String() != "[]" {
		t.Fatalf("expected empty array, got %s", resp.Body.String())
	}
}
