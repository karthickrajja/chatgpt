package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"time"

	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type user struct {
	ID        uint      `json:"id" gorm:"primary_key;auto_increment"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"type:DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

type users struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var DB *gorm.DB

func databaseconnection() {
	database, err := gorm.Open("mysql", "root:@Kar9600@tcp(localhost:3306)/AI_writer?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("couldn't connect with database")
	}
	err = database.AutoMigrate(&user{}).Error
	if err != nil {
		fmt.Println("Database Connection Issue")
		return
	}
	DB = database
}

var otp = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func generateotp(length int) string {
	otp_number := make([]byte, 6)
	n, err := io.ReadAtLeast(rand.Reader, otp_number, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < length; i++ {
		otp_number = append(otp_number, otp[i])
	}
	return string(otp_number)
}

func Createuser(c *gin.Context) {
	var inputs users
	if err := c.BindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	data := user{Email: inputs.Email, Password: inputs.Password}
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	match := emailRegex.MatchString(inputs.Email)
	if !match {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Please enter the valid Email ID"})
		return
	}

	var existing user
	if err := DB.Where("email = ?", inputs.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Already a user please try Login"})
		return
	}
	DB.Create(&data)
	c.JSON(http.StatusOK, gin.H{"Message": "success"})
}
func retrivebyid(c *gin.Context) {
	var data user
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Please provide a correct input"})
		return
	}
	if err := DB.Where("id = ?", data.ID).Find(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No data Find"})
		return
	}
	c.JSON(http.StatusOK, data)
}

func tokens() {

}

func forgotpassword(c *gin.Context) {
	var data user
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Please enter a mail Id"})
		return
	}
	if err := DB.Where("email = ?", data.Email).First(&data).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "No data Find"})
	}
	c.JSON(http.StatusOK, data.Password)

}

func fetchdatabyemail(c *gin.Context) {
	var data user
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "Please enter the valid details"})
		return
	}
	if err := DB.Where("email = ?", data.Email).First(&data).Error; err == nil {
		c.JSON(http.StatusOK, data)
		return
	}
}

func main() {
	router := gin.Default()
	databaseconnection()
	router.POST("/signup", Createuser)
	router.GET("/data", retrivebyid)
	router.GET("/data/bymail", fetchdatabyemail)
	router.GET("/forgotpassword", forgotpassword)
	router.Run(":8080")
}
