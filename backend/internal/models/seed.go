package models

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount > 0 {
		fmt.Println("🔁 Seed skipped: data already exists")
		return
	}

	// Create test user
	hashed, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	user := User{
		Email:    "test@example.com",
		Password: string(hashed),
	}
	db.Create(&user)

	// FirstOrCreate Ingredients
	salt := Ingredient{Name: "Salt"}
	db.FirstOrCreate(&salt, salt)

	flour := Ingredient{Name: "Flour"}
	db.FirstOrCreate(&flour, flour)

	// FirstOrCreate Tags
	vegan := Tag{Name: "Vegan"}
	db.FirstOrCreate(&vegan, vegan)

	breakfast := Tag{Name: "Breakfast"}
	db.FirstOrCreate(&breakfast, breakfast)

	// Recipe
	recipe := Recipe{
		Title:        "Simple Pancakes",
		Instructions: "Mix ingredients and cook in pan.",
		CreatedByID:  user.ID,
	}
	db.Create(&recipe)

	// Recipe <-> Ingredients
	db.Create(&[]RecipeIngredient{
		{RecipeID: recipe.ID, IngredientID: salt.ID, Amount: "1", Unit: "tsp"},
		{RecipeID: recipe.ID, IngredientID: flour.ID, Amount: "2", Unit: "cups"},
	})

	// Recipe <-> Tags
	db.Create(&[]RecipeTag{
		{RecipeID: recipe.ID, TagID: vegan.ID},
		{RecipeID: recipe.ID, TagID: breakfast.ID},
	})

	// Menu Entry
	db.Create(&MenuEntry{
		UserID:   user.ID,
		RecipeID: recipe.ID,
		Day:      "Monday",
		MealType: "breakfast",
	})

	fmt.Println("🌱 Seed complete")
}
