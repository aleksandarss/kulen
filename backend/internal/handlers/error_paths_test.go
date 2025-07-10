package handlers_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupDB() *gorm.DB {
	database, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.AutoMigrate(
		&models.User{},
		&models.Recipe{},
		&models.Ingredient{},
		&models.RecipeIngredient{},
		&models.Tag{},
		&models.RecipeTag{},
		&models.MenuEntry{},
		&models.MenuEntryExtra{},
		&models.RecipeStep{},
	)
	models.Seed(database)
	return database
}

func TestLoginBadPassword(t *testing.T) {
	dbConn := setupDB()
	db.DB = dbConn
	r := gin.Default()
	r.POST("/login", handlers.Login)
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(`{"email":"test@example.com","password":"wrong"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestRefreshTokenInvalid(t *testing.T) {
	dbConn := setupDB()
	db.DB = dbConn
	r := gin.Default()
	r.POST("/refresh", handlers.RefreshToken)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(`{"refresh_token":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetRecipeByIDNotFound(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.GET("/recipes/:id", handlers.GetRecipeByID(dbConn))
	req := httptest.NewRequest("GET", "/recipes/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateRecipeInvalidID(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.PUT("/recipes/:id", handlers.UpdateRecipe(dbConn))
	req := httptest.NewRequest("PUT", "/recipes/abc", bytes.NewBufferString(`{}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteRecipeInvalidID(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.DELETE("/recipes/:id", handlers.DeleteRecipe(dbConn))
	req := httptest.NewRequest("DELETE", "/recipes/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthMiddlewareMissing(t *testing.T) {
	r := gin.Default()
	r.GET("/protected", handlers.AuthMiddleware(), func(c *gin.Context) { c.String(200, "ok") })
	req := httptest.NewRequest("GET", "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
func TestCreateMenuEntryBadBody(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(1)); c.Next() })
	r.POST("/menu", handlers.CreateMenuEntry(dbConn))
	req := httptest.NewRequest("POST", "/menu", bytes.NewBufferString(`invalid`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteMenuEntryNotFound(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(1)); c.Next() })
	r.DELETE("/menu/:id", handlers.DeleteMenuEntry(dbConn))
	req := httptest.NewRequest("DELETE", "/menu/999", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateRecipeNotFound(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.PUT("/recipes/:id", handlers.UpdateRecipe(dbConn))
	req := httptest.NewRequest("PUT", "/recipes/999", bytes.NewBufferString(`{"title":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
func TestUpdateRecipeBadBody(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.PUT("/recipes/:id", handlers.UpdateRecipe(dbConn))
	req := httptest.NewRequest("PUT", "/recipes/1", bytes.NewBufferString(`bad`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateMenuEntryUnauthorized(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.POST("/menu", handlers.CreateMenuEntry(dbConn))
	req := httptest.NewRequest("POST", "/menu", bytes.NewBufferString(`{}`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestDeleteMenuEntryBadID(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.Use(func(c *gin.Context) { c.Set("userID", uint(1)); c.Next() })
	r.DELETE("/menu/:id", handlers.DeleteMenuEntry(dbConn))
	req := httptest.NewRequest("DELETE", "/menu/abc", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
func TestLoginInvalidBody(t *testing.T) {
	dbConn := setupDB()
	db.DB = dbConn
	r := gin.Default()
	r.POST("/login", handlers.Login)
	req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(`bad`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetMenuEntriesUnauthorized(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.GET("/menu", handlers.GetMenuEntries(dbConn))
	req := httptest.NewRequest("GET", "/menu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateRecipeTooManySteps(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.POST("/recipes", handlers.CreateRecipe(dbConn))
	var steps string
	for i := 0; i < 11; i++ {
		steps += `{"title":"a","text":"b"},`
	}
	body := `{"title":"t","instructions":"i","created_by_id":1,"steps":[` + steps + `]}`
	req := httptest.NewRequest("POST", "/recipes", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateRecipeInvalidJSON(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.POST("/recipes", handlers.CreateRecipe(dbConn))
	req := httptest.NewRequest("POST", "/recipes", bytes.NewBufferString(`bad`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
func TestRefreshTokenMissing(t *testing.T) {
	dbConn := setupDB()
	db.DB = dbConn
	r := gin.Default()
	r.POST("/refresh", handlers.RefreshToken)
	req := httptest.NewRequest("POST", "/refresh", bytes.NewBufferString(`bad`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateTagInvalidJSON(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.POST("/tags", handlers.CreateTag(dbConn))
	req := httptest.NewRequest("POST", "/tags", bytes.NewBufferString(`bad`))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
func TestGetIngredientsEmptyQuery(t *testing.T) {
	dbConn := setupDB()
	r := gin.Default()
	r.GET("/ingredients", handlers.GetIngredients(dbConn))
	req := httptest.NewRequest("GET", "/ingredients", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
func TestGetAllRecipesWithTag(t *testing.T) {
	dbConn := setupDB()
	r := setupTestRouter(dbConn)
	req := httptest.NewRequest("GET", "/recipes?tag=Vegan", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
