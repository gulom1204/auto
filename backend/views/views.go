package views

import (
	"autoparts/config"
	"autoparts/models"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
)
var jwtKey = []byte(os.Getenv("JWT_KEY"))

func GetHome(c *gin.Context) {
	var products []models.Product

	db, err := config.InitDB()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func Registration(c *gin.Context) {
	var user models.User
	// Инициализация базы данных
	db, err := config.InitDB()
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Проверка на уникальность email
	var existingUser models.User
	if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	// Если пользователь не найден, добавляем нового
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func Login(c *gin.Context) {
    var input models.Login
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ввод"})
        return
    }

    db, _ := config.InitDB()
    var user models.User
    if err := db.Where("email = ?", input.Email).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
        return
    }

    // Логируем введенные данные и хеш из базы данных для отладки
    log.Println("Проверка пароля:", input.Password)
    log.Println("Хеш пароля из базы данных:", user.Password)

    // Убираем пробелы с пароля перед сравнением
    inputPassword := strings.TrimSpace(input.Password)

    // Проверяем пароль
    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)); err != nil {
        log.Println("Ошибка сравнения паролей:", err)  // Логируем ошибку
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
        return
    }

    // Создание токена
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(time.Hour * 72).Unix(),
    })

    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать токен"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func CreateItems(c *gin.Context) {
	var products []models.Product
	db, err := config.InitDB()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&products); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&products).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"products": products})
}