package core

import (
	"errors"
	"regexp"
)

// ! Primary Port (customer_service.go)

type CustomerService interface {
	CreateCustomer(customer Customer) error
	GetCustomerById(customerId int) (*Customer, error)
	GetAllCustomer() ([]Customer, error)
	UpdateCustomer(customerId int, customer *Customer) (*Customer, error)
	DeleteCustomer(customerId int) error
	SearchCustomerById(customerId int) error
	ValidateName(customerName string) error
}

type customerServiceImpl struct {
	r CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return &customerServiceImpl{r: repo}
}

var (
	ErrInvalidFullName = errors.New("invalid full name")
	fullNameRegex      = `^[A-Z][a-z]+\s[A-Z][a-z]+(?:\s[A-Z][a-z]+)*$`
)

func ValidateFullName(fullName string) error {
	matched, err := regexp.MatchString(fullNameRegex, fullName)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidFullName
	}
	return nil
}

func (s *customerServiceImpl) CreateCustomer(customer Customer) error {
	// Business logic...
	// Register the custom validation function for 'Name'
	// validate.RegisterValidation("name", validateName)

	// if err := s.r.Validate(customer.Name); err != nil {
	// 	return err
	// }

	if customer.Age <= 0 {
		return errors.New("age must be positive")
	}

	if err := s.r.Save(customer); err != nil {
		return err
	}

	return nil
}

func (s *customerServiceImpl) GetCustomerById(customerId int) (*Customer, error) {
	// Business logic...
	if customerId < 0 {
		return &Customer{}, errors.New("customerId must be positive")
	}

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

func (s *customerServiceImpl) UpdateCustomer(customerId int, customer *Customer) (*Customer, error) {
	// Business logic...
	customer, err := s.r.Update(customerId, customer)

	if err != nil {
		return &Customer{}, err
	}

	return customer, nil
}

func (s *customerServiceImpl) DeleteCustomer(customerId int) error {
	// Business logic...
	if err := s.r.Delete(customerId); err != nil {
		return err
	}

	return nil
}

func (s *customerServiceImpl) SearchCustomerById(customerId int) error {
	// Business logic...
	if customerId < 0 {
		return errors.New("customerId must be positive")
	}

	if err := s.r.Search(customerId); err != nil {
		return err
	}

	return nil
}

var (
	ErrInvalidName = errors.New("invalid name")
	NameRegex      = `^[a-zA-Z\s]+$`
)

func (s *customerServiceImpl) ValidateName(customerName string) error {
	// Register the custom validation function for 'Name'
	matched, err := regexp.MatchString(NameRegex, customerName)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidName
	}
	return nil
}
