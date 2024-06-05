package core

import (
	"errors"
)

// ! Primary Port (customer_service.go)

type CustomerService interface {
	CreateCustomer(customer Customer) error
	GetCustomerById(customerId int) (*Customer, error)
	GetAllCustomer() ([]Customer, error)
}

type customerServiceImpl struct {
	r CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return &customerServiceImpl{r: repo}
}

func (s *customerServiceImpl) CreateCustomer(customer Customer) error {
	if customer.Age <= 0 {
		return errors.New("age must be positive")
	}
	// Business logic...
	if err := s.r.Save(customer); err != nil {
		return err
	}
	return nil
}

func (s *customerServiceImpl) GetCustomerById(customerId int) (*Customer, error) {
	if customerId < 0 {
		return &Customer{}, errors.New("customerId must be positive")
	}

	// Business logic...
	customer, err := s.r.Get(customerId)

	if err != nil {
		return &Customer{}, err
	}
	return customer, nil
}

func (s *customerServiceImpl) GetAllCustomer() ([]Customer, error) {
	// Business logic...
	var customers []Customer
	customers, err := s.r.GetAll()

	if err != nil {
		return []Customer{}, err
	}
	return customers, nil
}
