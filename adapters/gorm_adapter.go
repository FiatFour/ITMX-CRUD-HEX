package adapters

import (
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

	result := r.db.Table("customers").Find(&customers)

	if result.Error != nil {
		return []core.Customer{}, result.Error
	}
	// fmt.Println(customers)
	return customers, nil
}
