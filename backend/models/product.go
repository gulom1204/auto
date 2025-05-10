package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku" gorm:"unique"`
	Category    string  `json:"category"`
	Brand       string  `json:"brand"`
	Stock       int     `json:"stock"`
	ImageURL    string  `json:"image_url"`
} 
type User struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	LastName string `json:"lastname"`
	Email string `gorm:"unique" json:"email"`
	Phone string `gorm:"unique" json:"phone"`
	Password string `gorm:"unique" json:"password"`
	Role string `json:"role"`
}
type Login struct {
	Email string `json:"email"`
	Password string `json:"password"`
}