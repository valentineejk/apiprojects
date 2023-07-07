package main

import (
	"apiProject/models"
	"apiProject/storage"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create", r.CreateBook)
	api.Delete("/delete/:id", r.DeleteBook)
	api.Get("/book/:id", r.GetBook)
	api.Get("/books", r.GetBooks)

}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading env")
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DBNAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("error loading confg")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("could not migerate books")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := models.Book{}
	err := c.BodyParser(&book)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			fiber.Map{
				"status":  400,
				"message": "request failed",
				"data":    nil,
			})
		return err

	}
	err = r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  400,
			"message": "error book",
			"data":    nil,
		})
		return err
	}
	c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "book created",
		"data":    nil,
	})
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {

	bookModels := &[]models.Book{}

	err := r.DB.Find(bookModels).Error
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			fiber.Map{
				"status":  400,
				"message": "request failed",
				"data":    nil,
			})
		return err

	}

	c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
		"status":  200,
		"message": "book receisved",
		"data":    bookModels,
	})
	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {

	bookModel := models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  400,
			"message": "id cannot be empty",
			"data":    nil,
		})
		return nil
	}

	err := r.DB.Delete(bookModel, id)
	if err.Error != nil {

		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  400,
			"message": "failed to delete book",
			"data":    nil,
		})
		return nil
	}

	c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "book successfully deleted",
		"data":    nil,
	})
	return nil
}

func (r *Repository) GetBook(c *fiber.Ctx) error {
	id := c.Params("id")
	bookModel := &models.Book{}
	if id == "" {
		c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  400,
			"message": "id cannot be empty",
			"data":    nil,
		})
		return nil
	}
	fmt.Println(id)
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  400,
			"message": "error retrieving book",
			"data":    nil,
		})
		return nil
	}
	c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  200,
		"message": "book retrieved",
		"data":    bookModel,
	})
	return nil
}
