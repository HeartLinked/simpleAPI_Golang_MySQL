package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Book struct {
	//gorm.Model
	BookID   uint `gorm:"primaryKey;"`
	BookName string
	Author   string
	Year     uint8
	Number   int
}

func (b Book) TableName() string {
	return "book"
}

var db *gorm.DB
var err error

func init() {
	dsn := "root:LFYmemories0907@tcp(localhost:3306)/Library?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Book{})
}

func getData(c *gin.Context) {
	book_id := c.Query("book_id")
	book := Book{}
	tx := db.Where("book_id=?", book_id).First(&book)
	if tx.Error != nil {
		c.AbortWithStatus(404)
		panic(err.Error())
	}
	c.JSON(http.StatusOK, book)
}

func insertData(c *gin.Context) {
	var book Book
	err := c.Bind(&book)
	if err == nil {
		db.Create(&book)
		c.JSON(http.StatusOK, book)
	}
}

func updateData(c *gin.Context) {
	book_id := c.Param("book_id")
	action := c.Param("action")
	action = strings.Trim(action, "/")

	var book Book
	number, _ := strconv.Atoi(action)
	int_book_id, _ := strconv.Atoi(book_id)
	book.BookID = uint(int_book_id)
	db.Model(&book).Update("number", number)
	//c.JSON(200, gin.H{"number#" + number: "updated"})

}

func deleteData(c *gin.Context) {
	book_id := c.Params.ByName("book_id")
	var book Book
	d := db.Delete(&book, book_id)
	fmt.Println(d)
	c.JSON(200, gin.H{"book_id#" + book_id: "deleted"})
}

func main() {

	newbook := Book{BookID: 1001, BookName: "野草", Author: "luxun", Year: 12, Number: 2}

	if err := db.Create(&newbook).Error; err != nil {
		fmt.Println("插入失败", err)
		return
	}

	r := gin.Default()
	r.POST("/gorm", insertData)
	r.GET("/gorm", getData)
	r.PUT("/gorm/:book_id/*action", updateData)
	r.DELETE("/gorm/:book_id", deleteData)
	r.Run(":8080")
	// db.Create(&Product{})
}
