package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/db"
	"backend/internal/handlers"

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

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Simple health check route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API is up and running"})
	})

	// recipes
	router.GET("/recipes", handlers.GetAllRecipes(db.DB))
	router.GET("/recipes/:id", handlers.GetRecipeByID(db.DB))
	router.POST("/recipes", handlers.CreateRecipe(db.DB))
	router.PUT("/recipes/:id", handlers.UpdateRecipe(db.DB))
	router.DELETE("/recipes/:id", handlers.DeleteRecipe(db.DB))

	// menu
	router.GET("/menu", handlers.GetMenuEntries(db.DB))
	router.POST("/menu", handlers.CreateMenuEntry(db.DB))
	router.DELETE("/menu/:id", handlers.DeleteMenuEntry(db.DB))

	// shopping list
	router.GET("/shopping-list", handlers.GetShoppingList(db.DB))

	// ingredient
	router.GET("/ingredients", handlers.GetIngredients(db.DB))

	fmt.Printf("Starting server on port %s...\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
