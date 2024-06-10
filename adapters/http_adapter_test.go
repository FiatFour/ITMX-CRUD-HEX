package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/fiatfour/itmx-crud-hex/core"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCustomerService is a mock implementation of core.CustomerService
type MockCustomerService struct {
	mock.Mock
}

func (m *MockCustomerService) CreateCustomer(customer core.Customer) error {
	args := m.Called(customer)
	return args.Error(0)
}

func (m *MockCustomerService) GetCustomerById(customerId uint) (*core.Customer, error) {
	args := m.Called(customerId)
	return args.Get(0).(*core.Customer), args.Error(1)
}

func (m *MockCustomerService) GetAllCustomer() ([]core.Customer, error) {
	args := m.Called()
	return args.Get(0).([]core.Customer), args.Error(1)
}

func (m *MockCustomerService) UpdateCustomer(customerId uint, customer *core.Customer) (*core.Customer, error) {
	args := m.Called(customerId, customer)
	return args.Get(0).(*core.Customer), args.Error(1)
}

func (m *MockCustomerService) DeleteCustomer(customerId uint) error {
	args := m.Called(customerId)
	return args.Error(0)
}

func (m *MockCustomerService) SearchCustomerById(customerId uint) error {
	args := m.Called(customerId)
	return args.Error(0)
}

func (m *MockCustomerService) ValidateName(customerName string) error {
	args := m.Called(customerName)
	return args.Error(0)
}

// SetupTestApp initializes the Fiber app with the necessary routes and handlers for testing
func SetupTestApp(service core.CustomerService) *fiber.App {
	// initialize a new Fiber app
	app := fiber.New()

	// create a new handler with the provided service
	customerHandler := NewHttpCustomerHandler(service)

	// set up routes
	app.Post("/customers", customerHandler.CreateCustomerHandler)
	app.Get("/customers/:id", customerHandler.GetCustomerHandler)
	app.Get("/customers", customerHandler.GetAllCustomerHandler)
	app.Put("/customers/:id", customerHandler.UpdateCustomerHandler)
	app.Delete("/customers/:id", customerHandler.DeleteCustomerHandler)

	return app
}

func TestCreateCustomerHandler(t *testing.T) {
	// mock
	mockService := new(MockCustomerService)
	app := SetupTestApp(mockService)

	// Success case
	t.Run("successful create a customer ", func(t *testing.T) {
		// Mock service
		mockService.On("ValidateName", "Fiat").Return(nil)
		mockService.On("CreateCustomer", mock.AnythingOfType("core.Customer")).Return(nil)

		// create a new HTTP POST request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("POST", "/customers", bytes.NewBufferString(`{"name": "Fiat" ,"age": 24}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// check Error and Status
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	// Failure case
	t.Run("(fail) invalid request body", func(t *testing.T) {
		// clear mock
		mockService.ExpectedCalls = nil

		// create a new HTTP POST request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("POST", "/customers", bytes.NewBufferString(`{"name": 1 ,"age": "invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// check Error and Status
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	t.Run("(fail) name validation error", func(t *testing.T) {
		// clear mock
		mockService.ExpectedCalls = nil
		// Mock service
		mockService.On("ValidateName", "Invalid123").Return(core.ErrInvalidName)

		// create a new HTTP POST request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("POST", "/customers", bytes.NewBufferString(`{"name": "Invalid123", "age": 24}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// check Error and Status
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	t.Run("(fail) customer service error", func(t *testing.T) {
		// clear mock
		mockService.ExpectedCalls = nil
		// Mock service
		mockService.On("ValidateName", "Fiat").Return(nil)
		mockService.On("CreateCustomer", mock.AnythingOfType("core.Customer")).Return(errors.New("service error"))

		// create a new HTTP POST request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("POST", "/customers", bytes.NewBufferString(`{"name": "Fiat", "age": 24}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		// check Error and Status
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})
}

func TestGetCustomerHandler(t *testing.T) {
	// mock
	mockService := new(MockCustomerService)
	app := SetupTestApp(mockService)

	// Success case
	t.Run("successful get a customer", func(t *testing.T) {
		// setup Customer
		customerId := 1
		expectedCustomer := &core.Customer{ID: uint(customerId), Name: "Fiat", Age: 23}

		// mock service
		mockService.On("GetCustomerById", uint(customerId)).Return(expectedCustomer, nil)

		// create a new HTTP GET request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("GET", "/customers/"+strconv.Itoa(customerId), nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// decode JSON response from body and contain to Customer return Error and then check Value/Error
		var customer core.Customer
		err = json.NewDecoder(resp.Body).Decode(&customer)
		assert.NoError(t, err)
		assert.Equal(t, expectedCustomer, &customer)
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	// Failure case
	t.Run("(fail) invalid customer ID", func(t *testing.T) {
		// create a new HTTP GET request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("GET", "/customers/invalid-id", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("(fail) customer not found", func(t *testing.T) {
		// mock clear
		mockService.ExpectedCalls = nil
		// set up customerId
		customerId := 2
		// mock service
		mockService.On("GetCustomerById", uint(customerId)).Return(&core.Customer{}, errors.New("customer not found"))

		// create a new HTTP GET request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("GET", "/customers/"+strconv.Itoa(customerId), nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "customer not found", response["error"])
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})
}

func TestGetAllCustomerHandler(t *testing.T) {
	// mock
	mockService := new(MockCustomerService)
	app := SetupTestApp(mockService)

	// Success case
	t.Run("successful get all customers", func(t *testing.T) {
		// setup Customers
		expectedCustomers := []core.Customer{
			{ID: uint(1), Name: "Fiat", Age: uint(23)},
			{ID: uint(2), Name: "Anfat", Age: uint(40)},
		}

		// mock service
		mockService.On("GetAllCustomer").Return(expectedCustomers, nil)

		// create a new HTTP GET request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("GET", "/customers", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// decode JSON response from body and contain to Customers return Error and then check Value/Error
		var customers []core.Customer
		err = json.NewDecoder(resp.Body).Decode(&customers)
		assert.NoError(t, err)
		// compare expectedCustomers with Customers
		for index, customer := range customers {
			assert.Equal(t, expectedCustomers[index].ID, customer.ID)
			assert.Equal(t, expectedCustomers[index].Name, customer.Name)
			assert.Equal(t, expectedCustomers[index].Age, customer.Age)
		}
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	// Failure case
	t.Run("(fail) get all customer service error", func(t *testing.T) {
		// clear mock
		mockService.ExpectedCalls = nil
		// Mock service
		mockService.On("GetAllCustomer").Return([]core.Customer{}, errors.New("service error"))

		// create a new HTTP GET request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("GET", "/customers", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "service error", response["error"])
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})
}

func TestUpdateCustomerHandler(t *testing.T) {
	// mock
	mockService := new(MockCustomerService)
	app := SetupTestApp(mockService)

	// Success case
	t.Run("successful update a customer", func(t *testing.T) {
		// setup Customer
		customerId := uint(1)
		updatedCustomer := &core.Customer{ID: customerId, Name: "Updated Name", Age: uint(24)}

		// mock service
		mockService.On("ValidateName", "Updated Name").Return(nil)
		mockService.On("SearchCustomerById", customerId).Return(nil)
		mockService.On("UpdateCustomer", customerId, mock.AnythingOfType("*core.Customer")).Return(updatedCustomer, nil)

		// create a new HTTP PUT request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("PUT", "/customers/1", bytes.NewBufferString(`{"name": "Updated Name", "age": 24}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Updated successfully!", response["status"])
		assert.Equal(t, float64(1), response["id"])
		assert.Equal(t, "Updated Name", response["name"])
		assert.Equal(t, float64(24), response["age"])
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	// Failure case
	t.Run("(fail) invalid request body", func(t *testing.T) {
		// create a new HTTP PUT request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("PUT", "/customers/1", bytes.NewBufferString(`{"name": 1 ,"age": "invalid"}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("(fail) invalid request", func(t *testing.T) {
		// create a new HTTP PUT request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("PUT", "/customers/invalid", bytes.NewBufferString(`{"name": 1 ,"age": "invalid"}`))
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request", response["error"])
	})

	t.Run("(fail) validation error", func(t *testing.T) {
		// clear mock
		mockService.ExpectedCalls = nil
		// Mock service
		mockService.On("ValidateName", "Invalid Name!").Return(errors.New("invalid name"))

		// create a new HTTP PUT request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("PUT", "/customers/1", bytes.NewBufferString(`{"name": "Invalid Name!", "age": 24}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid name", response["error"])
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	t.Run("(fail) customer not found", func(t *testing.T) {
		// mock clear
		mockService.ExpectedCalls = nil
		// setup customerId
		customerId := uint(1)
		// mock service
		mockService.On("ValidateName", "Fiat").Return(nil)
		mockService.On("SearchCustomerById", customerId).Return(errors.New("customer not found"))

		// create a new HTTP PUT request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("PUT", "/customers/1", bytes.NewBufferString(`{"name": "Fiat", "age": 24}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		// read the entire response body will return Response with Error and Check Value/Error
		response, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "customer not found", string(response))
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	t.Run("(fail) internal server error", func(t *testing.T) {
		// mock clear
		mockService.ExpectedCalls = nil
		// setup customerId
		customerId := uint(1)

		// mock service
		mockService.On("ValidateName", "Fiat").Return(nil)
		mockService.On("SearchCustomerById", customerId).Return(nil)
		mockService.On("UpdateCustomer", customerId, mock.AnythingOfType("*core.Customer")).Return(&core.Customer{}, errors.New("service error"))

		// create a new HTTP PUT request set JSON format and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("PUT", "/customers/1", bytes.NewBufferString(`{"name": "Fiat", "age": 24}`))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "service error", response["error"])
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

}

func TestDeleteCustomerHandler(t *testing.T) {
	// mock
	mockService := new(MockCustomerService)
	app := SetupTestApp(mockService)

	// Success case
	t.Run("successful delete a customer", func(t *testing.T) {
		// setup customerId
		customerId := uint(1)

		// mock service
		mockService.On("SearchCustomerById", customerId).Return(nil)
		mockService.On("DeleteCustomer", customerId).Return(nil)

		// create a new HTTP Delete request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("DELETE", "/customers/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// read the entire response body will return Response with Error and Check Value/Error
		response, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "Deleted successfully!", string(response))
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	// Failure case
	t.Run("(fail) invalid request", func(t *testing.T) {
		// create a new HTTP Delete request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("DELETE", "/customers/invalid", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "invalid request", response["error"])
	})

	t.Run("(fail) customer not found", func(t *testing.T) {
		// mock clear
		mockService.ExpectedCalls = nil
		// setup customerId
		customerId := uint(1)

		// mock service
		mockService.On("SearchCustomerById", customerId).Return(errors.New("customer not found"))

		// create a new HTTP Delete request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("DELETE", "/customers/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

		// read the entire response body will return Response with Error and Check Value/Error
		response, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Equal(t, "customer not found", string(response))
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})

	t.Run("(fail) internal server error", func(t *testing.T) {
		// mock clear
		mockService.ExpectedCalls = nil
		// setup customerId
		customerId := uint(1)

		// mock service
		mockService.On("SearchCustomerById", customerId).Return(nil)
		mockService.On("DeleteCustomer", customerId).Return(errors.New("service error"))

		// create a new HTTP Delete request and send that will return value of Response(Status) with Error to check
		req := httptest.NewRequest("DELETE", "/customers/1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		// decode JSON response from body and contain to Response return Error and then check Value/Error
		var response map[string]string
		err = json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "service error", response["error"])
		// check all mocked it's work on expected
		mockService.AssertExpectations(t)
	})
}
