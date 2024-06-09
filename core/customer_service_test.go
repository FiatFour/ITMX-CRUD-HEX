package core

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of CustomerRepository
type mockCustomerRepo struct {
	saveFunc         func(customer Customer) error
	getFunc          func(customerId uint) (*Customer, error)                     // Port
	getAllFunc       func() ([]Customer, error)                                   // Port
	updateFunc       func(customerId uint, customer *Customer) (*Customer, error) // Port
	deleteFunc       func(customerId uint) error                                  // Port
	searchFunc       func(customerId uint) error                                  // Port
	validateNameFunc func(customerName string) error                              // Port
}

func (m *mockCustomerRepo) Save(customer Customer) error {
	return m.saveFunc(customer)
}

func (m *mockCustomerRepo) Get(customerId uint) (*Customer, error) {
	return m.getFunc(customerId)
}

func (m *mockCustomerRepo) GetAll() ([]Customer, error) {
	return m.getAllFunc()
}

func (m *mockCustomerRepo) Update(customerId uint, customer *Customer) (*Customer, error) {
	return m.updateFunc(customerId, customer)
}

func (m *mockCustomerRepo) Delete(customerId uint) error {
	return m.deleteFunc(customerId)
}

func (m *mockCustomerRepo) Search(customerId uint) error {
	return m.searchFunc(customerId)
}

func (m *mockCustomerRepo) Validate(customerName string) error {
	return m.validateNameFunc(customerName)
}

func TestCreateCustomer(t *testing.T) {
	// Success case
	t.Run("successful", func(t *testing.T) {
		repo := &mockCustomerRepo{
			saveFunc: func(customer Customer) error {
				// Simulate successful
				return nil
			},
		}
		service := NewCustomerService(repo)

		// Create a Customer in service and check Error
		err := service.CreateCustomer(Customer{Name: "Fiat", Age: uint(24)})
		assert.NoError(t, err)
	})

	// Failure case
	t.Run("(fail) age must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			saveFunc: func(customer Customer) error {
				// Simulate successful
				return nil
			},
		}
		service := NewCustomerService(repo)

		// Create a Customer in service and check Error
		err := service.CreateCustomer(Customer{Name: "Fiat", Age: uint(0)})
		assert.Error(t, err)
		assert.Equal(t, "age must more than 0", err.Error())
	})

	t.Run("(fail) database error", func(t *testing.T) {
		repo := &mockCustomerRepo{
			saveFunc: func(customer Customer) error {
				// Simulate Failure
				return errors.New("database error")
			},
		}
		service := NewCustomerService(repo)

		// Create a Customer in service and check Error
		err := service.CreateCustomer(Customer{Name: "Fiat", Age: uint(24)})
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestGetCustomerById(t *testing.T) {
	// Success case
	t.Run("successful", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getFunc: func(customerId uint) (*Customer, error) {
				// Simulate successful
				return &Customer{ID: uint(1), Name: "Fiat", Age: uint(24)}, nil
			},
		}
		service := NewCustomerService(repo)

		// get a Customer from service by Id and check Value/Error
		customer, err := service.GetCustomerById(uint(1))
		assert.Equal(t, uint(1), customer.ID)
		assert.Equal(t, "Fiat", customer.Name)
		assert.Equal(t, uint(24), customer.Age)
		assert.NoError(t, err)
	})

	// Failure case
	t.Run("(fail) customerId must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getFunc: func(customerId uint) (*Customer, error) {
				// Simulate successful
				return &Customer{ID: uint(1), Name: "Fiat", Age: uint(24)}, nil
			},
		}
		service := NewCustomerService(repo)

		// get a Customer from service by Id and check Value/Error
		customer, err := service.GetCustomerById(uint(0))
		assert.Error(t, err)
		assert.Equal(t, Customer{}, *customer)
		assert.Equal(t, "customerId must more than 0", err.Error())
	})

	t.Run("(fail) database error", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getFunc: func(customerId uint) (*Customer, error) {
				// Simulate Failure
				return &Customer{}, errors.New("database error")
			},
		}
		service := NewCustomerService(repo)

		// get a Customer from service by Id and check Value/Error
		customer, err := service.GetCustomerById(uint(1))
		assert.Error(t, err)
		assert.Equal(t, Customer{}, *customer)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestGetAllCustomer(t *testing.T) {
	// Success case
	t.Run("successful", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getAllFunc: func() ([]Customer, error) {
				// Simulate successful
				return []Customer{
					{ID: uint(1), Name: "Fiat", Age: uint(24)},
					{ID: uint(2), Name: "Anfat", Age: uint(40)},
				}, nil
			},
		}
		service := NewCustomerService(repo)

		expectedCustomers := []Customer{
			{ID: uint(1), Name: "Fiat", Age: uint(24)},
			{ID: uint(2), Name: "Anfat", Age: uint(40)},
		}

		// get all Customers from service by Id and and check Value/Error
		customers, err := service.GetAllCustomer()
		assert.NoError(t, err)

		// compare values
		assert.Len(t, customers, len(expectedCustomers))
		for index, customer := range customers {
			assert.Equal(t, expectedCustomers[index].ID, customer.ID)
			assert.Equal(t, expectedCustomers[index].Name, customer.Name)
			assert.Equal(t, expectedCustomers[index].Age, customer.Age)
		}
	})

	// Failure case
	t.Run("(fail) database error", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getAllFunc: func() ([]Customer, error) {
				// Simulate Failure
				return []Customer{}, errors.New("database error")
			},
		}
		service := NewCustomerService(repo)

		// get all Customers from service by Id and check Value/Error
		customers, err := service.GetAllCustomer()
		assert.Error(t, err)
		assert.Equal(t, []Customer{}, customers)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestUpdateCustomer(t *testing.T) {
	// Success case
	t.Run("successful", func(t *testing.T) {
		repo := &mockCustomerRepo{
			updateFunc: func(customerId uint, customer *Customer) (*Customer, error) {
				// Simulate successful
				return &Customer{ID: uint(1), Name: "Fiat", Age: uint(24)}, nil
			},
		}
		service := NewCustomerService(repo)

		// update a customer in service by Id with Customer and check Value/Error
		customer, err := service.UpdateCustomer(uint(1), &Customer{Name: "Fiat", Age: uint(24)})
		assert.NoError(t, err)
		assert.Equal(t, uint(1), customer.ID)
		assert.Equal(t, "Fiat", customer.Name)
		assert.Equal(t, uint(24), customer.Age)
	})

	// Fail case
	t.Run("(fail) age must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			updateFunc: func(customerId uint, customer *Customer) (*Customer, error) {
				// Simulate successful
				return &Customer{ID: uint(1), Name: "Fiat", Age: uint(24)}, nil
			},
		}
		service := NewCustomerService(repo)

		// update a customer in service by Id with Customer  and check Value/Error
		updatedCustomer, err := service.UpdateCustomer(uint(1), &Customer{Name: "Anfat", Age: uint(0)})
		assert.Error(t, err)
		assert.NotEqual(t, "Anfat", updatedCustomer.Name)
		assert.Equal(t, "age must more than 0", err.Error())
	})

	t.Run("(fail) database error", func(t *testing.T) {
		repo := &mockCustomerRepo{
			updateFunc: func(customerId uint, customer *Customer) (*Customer, error) {
				// Simulate failure
				return &Customer{}, errors.New("database error")
			},
		}
		service := NewCustomerService(repo)

		// update a customer in service by Id with Customer and check Value/Error
		updatedCustomer, err := service.UpdateCustomer(uint(1), &Customer{Name: "Anfat", Age: uint(40)})
		assert.Error(t, err)
		assert.Equal(t, &Customer{}, updatedCustomer)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestDeleteCustomer(t *testing.T) {
	// Success case
	t.Run("successful", func(t *testing.T) {
		repo := &mockCustomerRepo{
			deleteFunc: func(customerId uint) error {
				// Simulate successful
				return nil
			},
		}
		service := NewCustomerService(repo)

		// delete a customer in service by Id and check Error
		err := service.DeleteCustomer(uint(1))
		assert.NoError(t, err)
	})

	// Fail case
	t.Run("(fail) database error", func(t *testing.T) {
		repo := &mockCustomerRepo{
			deleteFunc: func(customerId uint) error {
				// Simulate failure
				return errors.New("database error")
			},
		}
		service := NewCustomerService(repo)

		// delete a customer in service by Id and check Error
		err := service.DeleteCustomer(uint(1))
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestSearchCustomerById(t *testing.T) {
	// Success case
	t.Run("successful", func(t *testing.T) {
		repo := &mockCustomerRepo{
			searchFunc: func(customerId uint) error {
				// Simulate successful
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.SearchCustomerById(uint(1))
		assert.NoError(t, err)
	})

	// Failure case
	t.Run("(fail) customerId must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			searchFunc: func(customerId uint) error {
				// Simulate successful
				return nil
			},
		}
		service := NewCustomerService(repo)

		// search a customer in service by Id and check Error
		err := service.SearchCustomerById(uint(0))
		assert.Error(t, err)
		assert.Equal(t, "customerId must more than 0", err.Error())
	})

	t.Run("(fail) database error", func(t *testing.T) {
		repo := &mockCustomerRepo{
			searchFunc: func(customerId uint) error {
				// Simulate failure
				return errors.New("database error")
			},
		}
		service := NewCustomerService(repo)

		// search a customer in service by Id and check Error
		err := service.SearchCustomerById(uint(1))
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
	})
}

func TestValidateName(t *testing.T) {
	// Success case
	t.Run("successful", func(t *testing.T) {
		repo := &mockCustomerRepo{
			validateNameFunc: func(customerName string) error {
				// Simulate successful
				return nil
			},
		}
		service := NewCustomerService(repo)

		// validate name and check Error
		err := service.ValidateName("Anfat Nilaingan")
		assert.NoError(t, err)
	})

	// Failure case
	t.Run("(fail) invalid name", func(t *testing.T) {
		repo := &mockCustomerRepo{
			validateNameFunc: func(customerName string) error {
				// Simulate successful
				return nil
			},
		}
		service := NewCustomerService(repo)

		// validate name and check Error
		err := service.ValidateName("Anfat !@#$%")
		assert.Error(t, err)
		assert.Equal(t, "invalid name", err.Error())
	})
}
