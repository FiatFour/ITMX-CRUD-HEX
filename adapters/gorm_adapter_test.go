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
	db.Migrator().CreateTable(&core.Customer{})
	return db
}

func TestGormCustomerRepository_Save(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successful save", func(t *testing.T) {
		// setup Customer
		customer := core.Customer{Name: "Fiat", Age: uint(24)}
		// Save() for insert a Customer in database and check Error
		err := repo.Save(customer)
		assert.NoError(t, err)

		// Check a row has inserted
		var count int64
		db.Model(&core.Customer{}).Where("name = ?", customer.Name).Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("(fail) name already exists", func(t *testing.T) {
		// setup Customer
		customer := core.Customer{Name: "Fiat", Age: uint(24)}
		// Save() for insert a Customer in database and check Error
		err := repo.Save(customer)
		assert.Error(t, err)
		assert.Equal(t, "name already exists", err.Error())
	})

	t.Run("(fail) database error on insert", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Save() for insert a Customer in database and check Error
		err := repo.Save(core.Customer{Name: "Fiat", Age: 24})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_Get(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successful get", func(t *testing.T) {
		// Save() for insert a Customer in database and check Error
		err := repo.Save(core.Customer{Name: "Fiat", Age: 24})
		assert.NoError(t, err)

		// Get() for get a Customer by Id from database and check Value/Error
		getCustomer, err := repo.Get(1)
		assert.NoError(t, err)
		assert.Equal(t, "Fiat", getCustomer.Name)
		assert.Equal(t, uint(24), getCustomer.Age)
	})

	t.Run("(fail) customer not found", func(t *testing.T) {
		// Get() for get a Customer by Id from database and check Value/Error
		customer, err := repo.Get(uint(999))
		assert.Error(t, err)
		assert.Equal(t, &core.Customer{}, customer)
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})

	t.Run("(fail) database error on get", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Get() for get a Customer by Id from database and check Value/Error
		customer, err := repo.Get(uint(999))
		assert.Error(t, err)
		assert.Equal(t, &core.Customer{}, customer)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_GetAll(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successful get all customers", func(t *testing.T) {
		// setup Customers
		expectedCustomers := []core.Customer{
			{Name: "Fiat", Age: uint(24)},
			{Name: "Anfat", Age: uint(40)},
		}

		// Save() loop for insert Customers and check Error
		for _, customer := range expectedCustomers {
			err := repo.Save(customer)
			assert.NoError(t, err)
		}

		// get all Customers from database and check Value/Error
		getCustomers, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, getCustomers, 2)
		assert.Equal(t, expectedCustomers[0].Name, getCustomers[0].Name)
		assert.Equal(t, expectedCustomers[1].Name, getCustomers[1].Name)
	})

	t.Run("(fail) database error on get all", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// get all Customers from database and check Value/Error
		_, err := repo.GetAll()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_Update(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	// setup Customers
	customers := []core.Customer{
		{Name: "Fiat", Age: uint(24)},
		{Name: "Anfat", Age: uint(40)},
		{Name: "Nilaingan", Age: uint(70)},
		{Name: "Hi", Age: uint(20)},
	}

	t.Run("successful update", func(t *testing.T) {
		// Save() for insert a Customer in database and check Error
		err := repo.Save(customers[0])
		assert.NoError(t, err)

		// Check a row has inserted
		var count int64
		db.Model(&core.Customer{}).Where("name = ?", customers[0].Name).Count(&count)
		assert.Equal(t, int64(1), count)

		// Update() for update a Customer by Id with Customer[1] in database and check Value/Error
		updatedCustomer, err := repo.Update(uint(1), &customers[1])
		assert.NoError(t, err)
		assert.Equal(t, uint(1), updatedCustomer.ID)
		assert.Equal(t, &customers[1].Name, &updatedCustomer.Name)
		assert.Equal(t, &customers[1].Age, &updatedCustomer.Age)
	})

	t.Run("(fail) name already exists error", func(t *testing.T) {
		// Save() for insert a Customer in database and check Error
		err := repo.Save(customers[2])
		assert.NoError(t, err)

		// Update() for update a Customer by Id with Customer[1] in database and check Value/Error
		updatedCustomer, err := repo.Update(uint(1), &customers[2])
		assert.Equal(t, &core.Customer{}, updatedCustomer)
		assert.Error(t, err)
		assert.Equal(t, "name already exists", err.Error())
	})

	t.Run("(fail) database error on update", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Update() for update a Customer by Id with Customer[3] in database and check Value/Error
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
		// Save() for insert a Customer in database and check Error
		err := repo.Save(core.Customer{Name: "Fiat", Age: uint(24)})
		assert.NoError(t, err)

		// Check a row has inserted
		var count int64
		db.Model(&core.Customer{}).Where("name = ?", "Fiat").Count(&count)
		assert.Equal(t, int64(1), count)

		// Delete() for delete a Customer by Id in database and check Error
		err = repo.Delete(uint(999))
		assert.NoError(t, err)
	})

	t.Run("(fail) database error on delete", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Delete() for delete a Customer by Id in database and check Error
		err := repo.Delete(uint(1))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}

func TestGormCustomerRepository_Search(t *testing.T) {
	db := setupTestDB()
	repo := NewGormCustomerRepository(db)

	t.Run("successful search", func(t *testing.T) {
		// Save() for insert a Customer in database and check Error
		err := repo.Save(core.Customer{Name: "Fiat", Age: uint(24)})
		assert.NoError(t, err)

		// Check a row has inserted
		var count int64
		db.Model(&core.Customer{}).Where("name = ?", "Fiat").Count(&count)
		assert.Equal(t, int64(1), count)

		// Search() for search a Customer by Id from database and check Error
		err = repo.Search(uint(1))
		assert.NoError(t, err)
	})

	t.Run("(fail) not found on search", func(t *testing.T) {
		// Search() for search a Customer by Id from database and check Error
		err := repo.Search(uint(999))
		assert.Error(t, err)
	})

	t.Run("(fail) database error on search", func(t *testing.T) {
		// Close the database to force an error
		sqlDB, _ := db.DB()
		sqlDB.Close()

		// Search() for search a Customer by Id from database and check Value/Error
		err := repo.Search(uint(1))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database is closed")
	})
}
