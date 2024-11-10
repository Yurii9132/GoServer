package main

import (
	"fmt"
	"log"

	"github.com/kataras/iris/v12"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID             uint   `gorm:"primaryKey"`
	Name           string `gorm:"size:100;not null"`
	Brand          string `gorm:"size:50;not null"`
	Description    string `gorm:"type:text"`
	Price          string `gorm:"type:money;not null"`
	Stock          int    `gorm:"not null"`
	Category       string `gorm:"size:50;default:electronics"`
	Specifications string `gorm:"type:jsonb"`
	ImageURL       string `gorm:"type:text"`
	Colour         string `gorm:"size:50;not null"`
}

var db *gorm.DB

func initDB() {
	var err error
	dsn := "user=me password=password dbname=electronic_devices host=localhost port=5432 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	fmt.Println("Connected to database")
	db.AutoMigrate(&Product{})
}

func main() {
	initDB()

	app := iris.New()
	app.RegisterView(iris.HTML("./views", ".html"))

	// Serve the styles directory as /static
	app.HandleDir("/static", "./styles")

	app.Get("/products", getProducts)

	app.Listen(":8080")
}

func getProducts(ctx iris.Context) {
	var products []Product
	if err := db.Find(&products).Error; err != nil {
		ctx.StatusCode(500)
		ctx.WriteString("Error getting products")
		return
	}
	// Render the HTML template with the products data
	ctx.ViewData("Products", products) // Passing data to the view
	ctx.View("products.html")          // Render the "products.html" template
}
