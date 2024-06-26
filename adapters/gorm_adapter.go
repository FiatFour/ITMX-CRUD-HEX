package adapters

import (
	"errors"

	"github.com/fiatfour/itmx-crud-hex/core"
	"gorm.io/gorm"
)

// * Secondary adapter (gorm_adapter.go)

type GormCustomerRepository struct {
	db *gorm.DB
}

func NewGormCustomerRepository(db *gorm.DB) core.CustomerRepository {
	return &GormCustomerRepository{db: db}
}

func (r *GormCustomerRepository) Save(customer core.Customer) error {
	var count int64

	// Check name is exists or not and check Error
	if err := r.db.Model(&core.Customer{}).Where("name = ?", customer.Name).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("name already exists")
	}

	// Insert Customer in database and check Error
	if err := r.db.Create(&customer).Error; err != nil {
		// Handle database errors
		return err
	}

	return nil
}

func (r *GormCustomerRepository) Get(customerId uint) (*core.Customer, error) {
	var customer core.Customer

	// Get a Customer from database and check Error
	if err := r.db.First(&customer, customerId).Error; err != nil {
		return &core.Customer{}, err
	}

	return &customer, nil
}

func (r *GormCustomerRepository) GetAll() ([]core.Customer, error) {
	var customers []core.Customer

	// Get all Customers from database and check Error
	if err := r.db.Find(&customers).Error; err != nil {
		return []core.Customer{}, err
	}

	return customers, nil
}

func (r *GormCustomerRepository) Update(customerId uint, customer *core.Customer) (*core.Customer, error) {
	var count int64
	// Check name is exists or not except the customerId for update and check Error
	if err := r.db.Model(&core.Customer{}).Where("id != ? AND name = ?", customer.ID, customer.Name).Count(&count).Error; err != nil {
		return &core.Customer{}, err
	}
	if count > 0 {
		return &core.Customer{}, errors.New("name already exists")
	}

	// Update a Customer in database and check Error
	if err := r.db.Model(&core.Customer{}).Where("id = ?", customerId).Updates(customer).Error; err != nil {
		return &core.Customer{}, err
	}
	customer.ID = uint(customerId)

	return customer, nil
}

func (r *GormCustomerRepository) Delete(customerId uint) error {
	// Delete a Customer in database from customerId and check Error
	if err := r.db.Where("id = ?", customerId).Delete(&core.Customer{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *GormCustomerRepository) Search(customerId uint) error {
	// Search a Customer in database from customerId and check Error
	if err := r.db.First(&core.Customer{}, customerId).Error; err != nil {
		return err
	}

	return nil
}
