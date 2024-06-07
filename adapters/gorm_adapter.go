package adapters

import (
	"errors"
	"fmt"

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

	if err := r.db.Model(&core.Customer{}).Where("name = ?", customer.Name).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("name already exists")
	}

	if err := r.db.Create(&customer).Error; err != nil {
		// Handle database errors
		return err
	}

	return nil
}

func (r *GormCustomerRepository) Get(customerId int) (*core.Customer, error) {
	var customer core.Customer

	if err := r.db.First(&customer, customerId).Error; err != nil {
		return &core.Customer{}, err
	}
	// fmt.Println(customer)

	return &customer, nil
}

func (r *GormCustomerRepository) GetAll() ([]core.Customer, error) {
	var customers []core.Customer

	if err := r.db.Find(&customers).Error; err != nil {
		return []core.Customer{}, err
	}
	// fmt.Println(customers)s

	return customers, nil
}

func (r *GormCustomerRepository) Update(customerId int, customer *core.Customer) (*core.Customer, error) {
	var count int64
	if err := r.db.Model(&core.Customer{}).Where("id != ? AND name = ?", customer.ID, customer.Name).Count(&count).Error; err != nil {
		return &core.Customer{}, err
	}
	if count > 0 {
		return &core.Customer{}, errors.New("name already exists")
	}

	if err := r.db.Model(&core.Customer{}).Where("id = ?", customerId).Updates(customer).Error; err != nil {
		return &core.Customer{}, err
	}
	customer.ID = uint(customerId)
	fmt.Println(customer)

	return customer, nil
}

func (r *GormCustomerRepository) Delete(customerId int) error {
	if err := r.db.Where("id = ?", customerId).Delete(&core.Customer{}).Error; err != nil {
		return err
	}
	// fmt.Println(customer)
	return nil
}

func (r *GormCustomerRepository) Search(customerId int) error {
	if err := r.db.First(&core.Customer{}, customerId).Error; err != nil {
		return err
	}

	return nil
}
