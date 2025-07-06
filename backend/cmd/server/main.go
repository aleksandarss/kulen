package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Init DB
	db.Init()

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())
	router.RedirectTrailingSlash = false

	api := router.Group("/api")

	// cors setup
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081", "https://kulen.mithrandir.calic.cloud"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Simple health check route
	api.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is up and running"})
	})

	// auth
	api.POST("/login", handlers.Login)
	api.POST("/refresh", handlers.RefreshToken)

	// recipes
	api.GET("/recipes", handlers.GetAllRecipes(db.DB))
	api.GET("/recipes/:id", handlers.GetRecipeByID(db.DB))
	api.POST("/recipes", handlers.CreateRecipe(db.DB))
	api.PUT("/recipes/:id", handlers.UpdateRecipe(db.DB))
	api.DELETE("/recipes/:id", handlers.DeleteRecipe(db.DB))

	// menu
	menu := api.Group("/menu")
	menu.Use(middleware.AuthRequired())
	{
		menu.GET("", handlers.GetMenuEntries(db.DB))  // handles /menu
		menu.GET("/", handlers.GetMenuEntries(db.DB)) // handles /menu/
		menu.POST("", handlers.CreateMenuEntry(db.DB))
		menu.POST("/", handlers.CreateMenuEntry(db.DB))
		menu.DELETE("/:id", handlers.DeleteMenuEntry(db.DB))
	}

	// shopping list
	api.GET("/shopping-list", handlers.GetShoppingList(db.DB))

	// ingredient
	api.GET("/ingredients", handlers.GetIngredients(db.DB))
	// tags
	api.GET("/tags", handlers.GetTags(db.DB))
	api.POST("/tags", handlers.CreateTag(db.DB))

	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
