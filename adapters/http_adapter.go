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

	// get a Customer from body(json) and check Error
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Validate name in service and check Error
	if err := h.service.ValidateName(customer.Name); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// call CreateCustomer() to pass agreement of Customer for create in service and check Error
	if err := h.service.CreateCustomer(customer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).SendString("Created successfully!")
}

func (h *HttpCustomerHandler) GetCustomerHandler(c *fiber.Ctx) error {
	// get Id and check Error
	customerId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// call GetCustomerById() to pass agreement of customerId for get a Customer in service and check Error
	customer, err := h.service.GetCustomerById(uint(customerId))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(customer)
}

func (h *HttpCustomerHandler) GetAllCustomerHandler(c *fiber.Ctx) error {
	var customers []core.Customer

	// call GetCustomerById() to pass agreement of customerId for get Customers in service and check Error
	customers, err := h.service.GetAllCustomer()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(customers)
}

func (h *HttpCustomerHandler) UpdateCustomerHandler(c *fiber.Ctx) error {
	var customer core.Customer

	// get Id and check error
	customerId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// get a Customer from body(json) and check Error
	if err := c.BodyParser(&customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// Validate name in service and check Error
	if err := h.service.ValidateName(customer.Name); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})

	}

	// call SearchCustomerById() to pass agreement of customerId for search a customer in service and check error
	if err = h.service.SearchCustomerById(uint(customerId)); err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	// call UpdateCustomer() to pass agreement of customerId with Customer for update a customer in service and get updatedCustomer with check error
	updatedCustomer, err := h.service.UpdateCustomer(uint(customerId), &customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "Updated successfully!",
		"id":     updatedCustomer.ID,
		"name":   updatedCustomer.Name,
		"age":    updatedCustomer.Age,
	})
}

func (h *HttpCustomerHandler) DeleteCustomerHandler(c *fiber.Ctx) error {
	// get Id and check error
	customerId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	// call SearchCustomerById() to pass agreement of customerId for search a customer in service and check error
	if err = h.service.SearchCustomerById(uint(customerId)); err != nil {
		return c.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	// call DeleteCustomer() to pass agreement of customerId for delete a customer in service and check error
	if err = h.service.DeleteCustomer(uint(customerId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).SendString("Deleted successfully!")
}
