package models

type Recipe struct {
	ID           uint   `gorm:"primaryKey"`
	Title        string `gorm:"not null"`
	Instructions string `gorm:"type:text"`
	CreatedByID  *uint  // Optional: who created it
	Ingredients  []RecipeIngredient
	Tags         []RecipeTag
	Steps        []RecipeStep `gorm:"constraint:OnDelete:CASCADE;"`
}

type RecipeStep struct {
	ID       uint   `gorm:"primaryKey"`
	RecipeID uint   `gorm:"not null"`
	Number   int    `gorm:"not null"`
	Title    string `gorm:"not null;default:'Step'"`
	Text     string `gorm:"type:text;not null"`

	Recipe Recipe `gorm:"constraint:OnDelete:CASCADE;"`
}

type Ingredient struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

type RecipeIngredient struct {
	ID           uint `gorm:"primaryKey"`
	RecipeID     uint
	IngredientID uint
	Amount       string
	Unit         string

	Recipe     Recipe
	Ingredient Ingredient
}

type Tag struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"uniqueIndex;not null"`
}

type RecipeTag struct {
	ID       uint `gorm:"primaryKey"`
	RecipeID uint
	TagID    uint

	Recipe Recipe
	Tag    Tag
}

type MenuEntry struct {
	ID       uint   `gorm:"primaryKey"`
	UserID   uint   `gorm:"not null;uniqueIndex:idx_user_day_meal"`
	RecipeID uint   `gorm:"not null"`
	Day      string `gorm:"not null;uniqueIndex:idx_user_day_meal"` // e.g., Monday-Sunday
	MealType string `gorm:"not null;uniqueIndex:idx_user_day_meal"` // e.g., breakfast, lunch, dinner

	Recipe Recipe
	User   User             `gorm:"constraint:OnDelete:CASCADE;"`
	Extras []MenuEntryExtra `gorm:"constraint:OnDelete:CASCADE;"`
}

type MenuEntryExtra struct {
	ID          uint   `gorm:"primaryKey"`
	MenuEntryID uint   `gorm:"not null;index"`
	Name        string `gorm:"not null"`

	MenuEntry MenuEntry `gorm:"constraint:OnDelete:CASCADE;"`
}

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"uniqueIndex;not null"`
	Password     string `gorm:"not null"`
	RefreshToken string `gorm:"type:text"` // stores latest refresh token
}
