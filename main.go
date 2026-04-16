package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name     string    `gorm:"not null" json:"name"`
	Email    string    `gorm:"unique;not null" json:"email"`
	Password string    `gorm:"not null" json:"-"` // Does not appear in JSON response
}


func main() {
	// Harusnya menggunakan .env but for demonstration only
	dsn := "host=localhost user=postgres password=postgres dbname=go_workshop port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Karena user menggunakan uuid
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.AutoMigrate(&User{})
	fmt.Println("Database connected and migrated successfully!")

	r := gin.Default()


	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.POST("/users", func(c *gin.Context) {
		var input struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required,email"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		newUser := User{
			Name:     input.Name,
			Email:    input.Email,
			Password: input.Password, // Di dalam praktek, jangan lupa untuk menghash password terlebih dahulu!
		}

		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user", "details": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "data": newUser})
	})

	r.GET("/users", func(c *gin.Context) {
		var users []User

		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": users})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		var user User
		if err := db.First(&user, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": user,
		})
	})



	r.DELETE("/users/:id", func(c *gin.Context) {
		id := c.Param("id")

		if err := db.Delete(&User{}, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
	})

	fmt.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}