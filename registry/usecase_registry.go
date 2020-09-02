package registry

import (
	"sync"

	//"github.com/Lionparcel/cabeen-sabre/src/usecase"
	"github.com/jinzhu/gorm"
)

// UsecaseRegistry ...
type UsecaseRegistry interface {
	// User() usecase.User
	// Flight() usecase.Flight
}

type usecaseRegistry struct {
	repo RepositoryRegistry
}

// NewUsecaseRegistry ...
func NewUsecaseRegistry(db *gorm.DB) (r UsecaseRegistry) {
	var ucRegistry usecaseRegistry
	var once sync.Once

	once.Do(func() {
		repoReg := NewRepoRegistry(db)
		ucRegistry = usecaseRegistry{repo: repoReg}
	})
	return &ucRegistry
}

// User ...
// func (u *usecaseRegistry) User() (ucUser usecase.User) {
// 	loadOnce.Do(func() {
// 		ucUser = usecase.NewUserUc(u.repo.User())
// 	})
// 	return ucUser
// }

// func (u *usecaseRegistry) Flight() usecase.Flight {
// 	var once sync.Once
// 	var ucFlight usecase.Flight
// 	once.Do(func() {
// 		ucFlight = usecase.NewFlightUc(u.repo.SessionToken())
// 	})
// 	return ucFlight
// }
