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

func TestGetIngredients(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Ingredient{})
	db.Create(&models.Ingredient{Name: "Sugar"})

	r := gin.Default()
	r.GET("/ingredients", handlers.GetIngredients(db))

	req, _ := http.NewRequest("GET", "/ingredients?query=sug", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
