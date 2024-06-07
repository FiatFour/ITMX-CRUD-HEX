package core

//* Secondary Port (customer_repository.go)

type CustomerRepository interface { // Spec
	Save(customer Customer) error                                  // Port
	Get(customerId uint) (*Customer, error)                        // Port
	GetAll() ([]Customer, error)                                   // Port
	Update(customerId uint, customer *Customer) (*Customer, error) // Port
	Delete(customerId uint) error                                  // Port
	Search(customerId uint) error                                  // Port
}
