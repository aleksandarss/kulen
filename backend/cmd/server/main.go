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

	// cors setup
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Simple health check route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is up and running"})
	})

	// auth
	router.POST("/login", handlers.Login)
	router.POST("/refresh", handlers.RefreshToken)

	// recipes
	router.GET("/recipes", handlers.GetAllRecipes(db.DB))
	router.GET("/recipes/:id", handlers.GetRecipeByID(db.DB))
	router.POST("/recipes", handlers.CreateRecipe(db.DB))
	router.PUT("/recipes/:id", handlers.UpdateRecipe(db.DB))
	router.DELETE("/recipes/:id", handlers.DeleteRecipe(db.DB))

	// menu
	menu := router.Group("/menu")
	menu.Use(middleware.AuthRequired())
	{
		menu.GET("", handlers.GetMenuEntries(db.DB))  // handles /menu
		menu.GET("/", handlers.GetMenuEntries(db.DB)) // handles /menu/
		menu.POST("", handlers.CreateMenuEntry(db.DB))
		menu.POST("/", handlers.CreateMenuEntry(db.DB))
		menu.DELETE("/:id", handlers.DeleteMenuEntry(db.DB))
	}

	// shopping list
	router.GET("/shopping-list", handlers.GetShoppingList(db.DB))

	// ingredient
	router.GET("/ingredients", handlers.GetIngredients(db.DB))

	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
