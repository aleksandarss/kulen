package handlers_test

import (
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/auth"
	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAuthEnv() *gin.Engine {
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.AutoMigrate(&models.User{})
	password, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	database.Create(&models.User{Email: "test@example.com", Password: string(password)})
	db.DB = database
	r := gin.Default()
	r.POST("/login", handlers.Login)
	r.POST("/refresh", handlers.RefreshToken)
	r.GET("/protected", handlers.AuthMiddleware(), func(c *gin.Context) {
		c.String(200, "ok")
	})
	return r
}

func TestLoginAndRefresh(t *testing.T) {
	r := setupAuthEnv()

	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"email":"test@example.com","password":"secret"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	refresh := resp["refresh_token"]
	if refresh == "" {
		t.Fatal("no refresh token returned")
	}

	req2 := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(`{"refresh_token":"`+refresh+`"}`))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w2.Code)
	}
}

func TestAuthMiddleware(t *testing.T) {
	r := setupAuthEnv()
	token, _ := auth.GenerateJWT(1)
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
