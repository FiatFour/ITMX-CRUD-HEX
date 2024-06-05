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

	if err := h.service.CreateCustomer(customer); err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}

func (h *HttpCustomerHandler) GetCustomerHandler(c *fiber.Ctx) error {
	// var customer core.Customer
	// customer := new(core.Customer)

	customerId, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	customer, err := h.service.GetCustomerById(customerId)
	if err != nil {
		// Return an appropriate error message and status code
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}
