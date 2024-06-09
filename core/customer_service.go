package core

import (
	"errors"
	"regexp"
)

// ! Primary Port (customer_service.go)
type CustomerService interface {
	CreateCustomer(customer Customer) error
	GetCustomerById(customerId uint) (*Customer, error)
	GetAllCustomer() ([]Customer, error)
	UpdateCustomer(customerId uint, customer *Customer) (*Customer, error)
	DeleteCustomer(customerId uint) error
	SearchCustomerById(customerId uint) error
	ValidateName(customerName string) error
}

// Implement CustomerRepository
type customerServiceImpl struct {
	r CustomerRepository
}

func NewCustomerService(repo CustomerRepository) CustomerService {
	return &customerServiceImpl{r: repo}
}

func (s *customerServiceImpl) CreateCustomer(customer Customer) error {
	// Business logic...
	// Check Age
	if customer.Age == 0 {
		return errors.New("age must more than 0")
	}

	// call Save() to pass agreement value of Customer for insert in gorm adapter
	if err := s.r.Save(customer); err != nil {
		return err
	}

	return nil
}

func (s *customerServiceImpl) GetCustomerById(customerId uint) (*Customer, error) {
	// Business logic...
	// Check customerId
	if customerId == 0 {
		return &Customer{}, errors.New("customerId must more than 0")
	}

	// call Get() to pass agreement customerId for get a Customer from gorm adapter
	customer, err := s.r.Get(customerId)
	if err != nil {
		return &Customer{}, err
	}

	return customer, nil
}

func (s *customerServiceImpl) GetAllCustomer() ([]Customer, error) {
	// Business logic...
	var customers []Customer

	// call GetAll() for get all Customers from gorm adapter
	customers, err := s.r.GetAll()
	if err != nil {
		return []Customer{}, err
	}

	return customers, nil
}

func (s *customerServiceImpl) UpdateCustomer(customerId uint, customer *Customer) (*Customer, error) {
	// Business logic...
	// Check Age
	if customer.Age == 0 {
		return &Customer{}, errors.New("age must more than 0")
	}

	// call Update() to pass agreement customerId and Customer for update a customer in gorm adapter and return value updated
	customer, err := s.r.Update(customerId, customer)

	if err != nil {
		return &Customer{}, err
	}

	return customer, nil
}

func (s *customerServiceImpl) DeleteCustomer(customerId uint) error {
	// Business logic...
	// call Delete() to pass agreement customerId for delete a customer in gorm adapter
	if err := s.r.Delete(customerId); err != nil {
		return err
	}

	return nil
}

func (s *customerServiceImpl) SearchCustomerById(customerId uint) error {
	// Business logic...
	// Check customerId
	if customerId == 0 {
		return errors.New("customerId must more than 0")
	}

	// call Search() to pass agreement customerId for search a customer in gorm adapter
	if err := s.r.Search(customerId); err != nil {
		return err
	}

	return nil
}

// define for validate name for checks if the value contains only alphabets and spaces
var (
	ErrInvalidName = errors.New("invalid name")
	NameRegex      = `^[a-zA-Z\s]+$`
)

func (s *customerServiceImpl) ValidateName(customerName string) error {
	// Validate name and check
	matched, err := regexp.MatchString(NameRegex, customerName)
	if err != nil {
		return err
	}
	if !matched {
		return ErrInvalidName
	}
	return nil
}
