package adapters

import (
	"strconv"

	"github.com/fiatfour/itmx-crud-hex/core"
	"github.com/gofiber/fiber/v2"
)

// ! Primary adapter (http_adapter.go)

type HttpCustomerHandler struct {
	service core.CustomerService
}

func NewHttpCustomerHandler(service core.CustomerService) *HttpCustomerHandler {
	return &HttpCustomerHandler{service: service}
}

func (h *HttpCustomerHandler) CreateCustomerHandler(c *fiber.Ctx) error {
	var customer core.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.service.ValidateName(customer.Name); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	if err := h.service.CreateCustomer(customer); err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// return c.Status(fiber.StatusCreated).JSON(customer)
	return c.Status(fiber.StatusCreated).SendString("Created successfully!")
}

func (h *HttpCustomerHandler) GetCustomerHandler(c *fiber.Ctx) error {
	customerId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	customer, err := h.service.GetCustomerById(uint(customerId))
	if err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}

func (h *HttpCustomerHandler) GetAllCustomerHandler(c *fiber.Ctx) error {
	var customers []core.Customer

	customers, err := h.service.GetAllCustomer()
	if err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(customers)
}

func (h *HttpCustomerHandler) UpdateCustomerHandler(c *fiber.Ctx) error {
	var customer core.Customer
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	customerId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := h.service.ValidateName(customer.Name); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	if err = h.service.SearchCustomerById(uint(customerId)); err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	updatedCustomer, err := h.service.UpdateCustomer(uint(customerId), &customer)
	if err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// return c.Status(fiber.StatusCreated).JSON(updatedCustomer)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "Updated successfully!",
		"id":     updatedCustomer.ID,
		"name":   updatedCustomer.Name,
		"age":    updatedCustomer.Age,
	})
}

func (h *HttpCustomerHandler) DeleteCustomerHandler(c *fiber.Ctx) error {
	customerId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err = h.service.SearchCustomerById(uint(customerId)); err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	if err = h.service.DeleteCustomer(uint(customerId)); err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString("Deleted successfully!")
}
