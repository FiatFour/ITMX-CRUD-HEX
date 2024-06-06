package core

type Customer struct {
	ID   uint `gorm:"primaryKey"`
	Name string
	Age  uint
}
