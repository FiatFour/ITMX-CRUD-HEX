package core

import (
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
	t.Run("success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			saveFunc: func(customer Customer) error {
				// Simulate successful save
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.CreateCustomer(Customer{Name: "Fiat", Age: 100})
		assert.NoError(t, err)
	})

	// Failure case: age must more than 0
	t.Run("age must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			saveFunc: func(customer Customer) error {
				// This won't be called due to validation
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.CreateCustomer(Customer{Name: "Fiat", Age: 0})
		assert.Error(t, err)
		assert.Equal(t, "age must more than 0", err.Error())
	})
}

func TestGetCustomerById(t *testing.T) {
	// Success case
	t.Run("success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getFunc: func(customerId uint) (*Customer, error) {
				// Simulate successful save
				return &Customer{ID: 1, Name: "Fiat", Age: 24}, nil
			},
		}
		service := NewCustomerService(repo)

		customer, err := service.GetCustomerById(1)
		assert.Equal(t, uint(1), customer.ID)
		assert.Equal(t, "Fiat", customer.Name)
		assert.Equal(t, uint(24), customer.Age)
		assert.NoError(t, err)
	})

	// Failure case: customerId must more than 0
	t.Run("customerId must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getFunc: func(customerId uint) (*Customer, error) {
				// Simulate successful save
				return &Customer{ID: 1, Name: "Fiat", Age: 24}, nil
			},
		}
		service := NewCustomerService(repo)

		customer, err := service.GetCustomerById(0)
		assert.Error(t, err)
		assert.Equal(t, Customer{}, *customer)
		assert.Equal(t, "customerId must more than 0", err.Error())
	})
}

func TestGetAllCustomer(t *testing.T) {
	// Success case
	t.Run("success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			getAllFunc: func() ([]Customer, error) {
				// Simulate successful save
				return []Customer{
					{ID: 1, Name: "Fiat", Age: 24},
					{ID: 2, Name: "Anfat", Age: 40},
				}, nil
			},
		}
		service := NewCustomerService(repo)

		expectedCustomers := []Customer{
			{ID: 1, Name: "Fiat", Age: 24},
			{ID: 2, Name: "Anfat", Age: 40},
		}

		customers, err := service.GetAllCustomer()

		assert.NoError(t, err)
		assert.Len(t, customers, len(expectedCustomers))

		for index, customer := range customers {
			assert.Equal(t, expectedCustomers[index].ID, customer.ID)
			assert.Equal(t, expectedCustomers[index].Name, customer.Name)
			assert.Equal(t, expectedCustomers[index].Age, customer.Age)
		}
	})
}

func TestUpdateCustomer(t *testing.T) {
	// Success case
	t.Run("success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			updateFunc: func(customerId uint, customer *Customer) (*Customer, error) {
				// Simulate successful save
				return &Customer{ID: 1, Name: "Fiat", Age: 24}, nil
			},
		}
		service := NewCustomerService(repo)

		customer, err := service.UpdateCustomer(uint(1), &Customer{Name: "Fiat", Age: 24})
		assert.NoError(t, err)

		assert.Equal(t, uint(1), customer.ID)
		assert.Equal(t, "Fiat", customer.Name)
		assert.Equal(t, uint(24), customer.Age)
	})

	// Fail case
	t.Run("age must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			updateFunc: func(customerId uint, customer *Customer) (*Customer, error) {
				// Simulate successful save
				return &Customer{ID: 1, Name: "Fiat", Age: 24}, nil
			},
		}
		service := NewCustomerService(repo)

		customer, err := service.UpdateCustomer(1, &Customer{Name: "Anfat", Age: 0})
		assert.Error(t, err)
		assert.NotEqual(t, "Anfat", customer.Name)
		assert.Equal(t, "age must more than 0", err.Error())
	})
}

func TestDeleteCustomer(t *testing.T) {
	// Success case
	t.Run("success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			deleteFunc: func(customerId uint) error {
				// Simulate successful save
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.DeleteCustomer(1)
		assert.NoError(t, err)
	})
}

func TestSearchCustomerById(t *testing.T) {
	// Success case
	t.Run("success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			searchFunc: func(customerId uint) error {
				// Simulate successful save
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.SearchCustomerById(1)
		assert.NoError(t, err)
	})

	// Failure case: customerId must more than 0
	t.Run("customerId must more than 0", func(t *testing.T) {
		repo := &mockCustomerRepo{
			searchFunc: func(customerId uint) error {
				// Simulate successful save
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.SearchCustomerById(0)
		assert.Error(t, err)
		assert.Equal(t, "customerId must more than 0", err.Error())
	})
}

func TestValidateName(t *testing.T) {
	// Success case
	t.Run("success", func(t *testing.T) {
		repo := &mockCustomerRepo{
			validateNameFunc: func(customerName string) error {
				// Simulate successful save
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.ValidateName("Anfat Nilaingan")
		assert.NoError(t, err)
	})

	// Failure case: invalid name
	t.Run("invalid name", func(t *testing.T) {
		repo := &mockCustomerRepo{
			validateNameFunc: func(customerName string) error {
				// Simulate successful save
				return nil
			},
		}
		service := NewCustomerService(repo)

		err := service.ValidateName("Anfat !@#$%")
		assert.Error(t, err)
		assert.Equal(t, "invalid name", err.Error())
	})
}
