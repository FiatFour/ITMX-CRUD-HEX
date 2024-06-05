package core

//* Secondary Port (customer_repository.go)

type CustomerRepository interface { // Spec
	Save(customer Customer) error          // Port
	Get(customerId int) (*Customer, error) // Port
	GetAll() ([]Customer, error)           // Port
}
