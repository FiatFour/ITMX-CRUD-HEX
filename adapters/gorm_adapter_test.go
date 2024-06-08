package adapters

import (
	"fmt"
	"testing"

	"github.com/fiatfour/itmx-crud-hex/core"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Setup an in-memory SQLite database with GORM
func setupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	db.AutoMigrate(&core.Customer{})
	return db
}

func TestGormCustomerRepository_Save(t *testing.T) {
	db := setupTestDB()
	// defer db.Rollback()

	repo := NewGormCustomerRepository(db)

	t.Run("successful save", func(t *testing.T) {
		customer := core.Customer{Name: "Fiat", Age: 24}
		err := repo.Save(customer)
		assert.NoError(t, err)

		var count int64
		db.Model(&core.Customer{}).Where("name = ?", customer.Name).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("name already exists error", func(t *testing.T) {
		customer := core.Customer{Name: "Fiat", Age: 24}
		err := repo.Save(customer)
		assert.Error(t, err)
		assert.Equal(t, "name already exists", err.Error())
	})

	t.Run("database error on insert", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		customer := core.Customer{Name: "Fiat", Age: 24}
		err := repo.Save(customer)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_Get(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successful save", func(t *testing.T) {
		customer := core.Customer{Name: "Fiat", Age: 24}
		err := repo.Save(customer)
		assert.NoError(t, err)

		var count int64
		db.Model(&core.Customer{}).Where("name = ?", customer.Name).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("customer not found", func(t *testing.T) {
		customer, err := repo.Get(999)
		assert.Error(t, err)
		assert.Equal(t, &core.Customer{}, customer)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("database error on insert", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		customer, err := repo.Get(999)
		assert.Error(t, err)
		assert.Equal(t, &core.Customer{}, customer)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_GetAll(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successfully get all customers", func(t *testing.T) {
		// Setup: Create some customers in the database
		customers := []core.Customer{
			{Name: "Fiat", Age: 24},
			{Name: "Anfat", Age: 40},
		}
		for _, customer := range customers {
			err := db.Create(&customer).Error
			assert.NoError(t, err)
		}

		// Act: get all customers
		getCustomers, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, getCustomers, 2)
		assert.Equal(t, customers[0].Name, getCustomers[0].Name)
		assert.Equal(t, customers[1].Name, getCustomers[1].Name)
	})

	t.Run("database error on get all", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		_, err := repo.GetAll()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_Update(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	customers := []core.Customer{
		{Name: "Fiat", Age: 24},
		{Name: "Anfat", Age: 40},
		{Name: "Nilaingan", Age: 70},
		{Name: "Hi", Age: 20},
	}

	t.Run("successful update", func(t *testing.T) {
		err := repo.Save(customers[0])
		assert.NoError(t, err)

		var count int64
		db.Model(&core.Customer{}).Where("name = ?", customers[0].Name).Count(&count)
		assert.Equal(t, int64(1), count)

		updatedCustomer, err := repo.Update(uint(1), &customers[1])
		assert.NoError(t, err)
		assert.Equal(t, uint(1), updatedCustomer.ID)
		assert.Equal(t, &customers[1].Name, &updatedCustomer.Name)
		assert.Equal(t, &customers[1].Age, &updatedCustomer.Age)
	})

	t.Run("name already exists error", func(t *testing.T) {
		err := repo.Save(customers[2])
		assert.NoError(t, err)
		updatedCustomer, err := repo.Update(uint(1), &customers[2])
		assert.Equal(t, &core.Customer{}, updatedCustomer)
		assert.Error(t, err)
		assert.Equal(t, "name already exists", err.Error())
	})

	t.Run("database error on update", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		updatedCustomer, err := repo.Update(uint(1), &customers[3])
		assert.Error(t, err)
		assert.Equal(t, &core.Customer{}, updatedCustomer)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_Delete(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successful delete", func(t *testing.T) {
		err := repo.Save(core.Customer{Name: "Fiat", Age: 24})
		assert.NoError(t, err)

		var count int64
		db.Model(&core.Customer{}).Where("name = ?", "Fiat").Count(&count)
		assert.Equal(t, int64(1), count)

		err = repo.Delete(uint(999))
		assert.NoError(t, err)
	})

	t.Run("delete when don't have value of that row", func(t *testing.T) {
		err := repo.Delete(uint(999))
		assert.NoError(t, err)
	})

	t.Run("database error on delete", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.Delete(uint(1))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_Search(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successful delete", func(t *testing.T) {
		err := repo.Save(core.Customer{Name: "Fiat", Age: 24})
		assert.NoError(t, err)

		var count int64
		db.Model(&core.Customer{}).Where("name = ?", "Fiat").Count(&count)
		assert.Equal(t, int64(1), count)

		err = repo.Search(uint(1))
		assert.NoError(t, err)
	})

	t.Run("not found for search", func(t *testing.T) {
		err := repo.Search(uint(999))
		assert.Error(t, err)
	})

	t.Run("database error on search", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		err := repo.Search(uint(1))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}
