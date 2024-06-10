package main

import (
	"github.com/fiatfour/itmx-crud-hex/adapters"
	"github.com/fiatfour/itmx-crud-hex/core"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// Initialize a new instance of a Fiber application
	app := fiber.New()

	// Initialize the database connection
	db, err := gorm.Open(sqlite.Open("customers.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema and insert rows of Customer
	db.Migrator().CreateTable(&core.Customer{})
	db.Create(&core.Customer{Name: "Fiat", Age: 24})
	db.Create(&core.Customer{Name: "Anfat Nilaingan", Age: 40})

	// Set up the core service and adapters
	customerRepo := adapters.NewGormCustomerRepository(db)
	customerService := core.NewCustomerService(customerRepo)
	customerHandler := adapters.NewHttpCustomerHandler(customerService)

	// Define routes
	app.Post("/customers", customerHandler.CreateCustomerHandler)
	app.Get("/customers/:id", customerHandler.GetCustomerHandler)
	app.Get("/customers", customerHandler.GetAllCustomerHandler)
	app.Put("/customers/:id", customerHandler.UpdateCustomerHandler)
	app.Delete("/customers/:id", customerHandler.DeleteCustomerHandler)

	// Start the server
	app.Listen("localhost:8080")
}
