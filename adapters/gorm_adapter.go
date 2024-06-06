package adapters

import (
	"fmt"

	"github.com/fiatfour/itmx-crud-hex/core"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// * Secondary adapter (gorm_adapter.go)

type GormCustomerRepository struct {
	db *gorm.DB
}

func NewGormCustomerRepository(db *gorm.DB) core.CustomerRepository {
	return &GormCustomerRepository{db: db}
}

func (r *GormCustomerRepository) Save(customer core.Customer) error {
	if result := r.db.Create(&customer); result.Error != nil {
		// Handle database errors
		return result.Error
	}
	return nil
}

func (r *GormCustomerRepository) Get(customerId int) (*core.Customer, error) {
	var customer core.Customer

	result := r.db.First(&customer, customerId)

	if result.Error != nil {
		return &core.Customer{}, result.Error
	}
	// fmt.Println(customer)
	return &customer, nil
}

func (r *GormCustomerRepository) GetAll() ([]core.Customer, error) {
	var customers []core.Customer

	result := r.db.Find(&customers)

	if result.Error != nil {
		return []core.Customer{}, result.Error
	}
	// fmt.Println(customers)s
	return customers, nil
}

func (r *GormCustomerRepository) Update(customerId int, customer *core.Customer) (*core.Customer, error) {

	// result := r.db.Model(&core.Customer{}).Where("id = ?", customerId).Updates(customer)
	result := r.db.Model(&core.Customer{}).Clauses(clause.Returning{}).Where("id = ?", customerId).Updates(customer)

	if result.Error != nil {
		return &core.Customer{}, result.Error
	}
	customer.ID = uint(customerId)
	fmt.Println(customer)
	return customer, nil
}

func (r *GormCustomerRepository) Delete(customerId int) error {
	result := r.db.Where("id = ?", customerId).Delete(&core.Customer{})

	if result.Error != nil {
		return result.Error
	}
	// fmt.Println(customer)
	return nil
}
