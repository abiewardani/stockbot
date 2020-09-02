package registry

import (
	"sync"

	"github.com/jinzhu/gorm"
)

// RepositoryRegistry ...
type RepositoryRegistry interface {
	// User() repository.User
	// SessionToken() repository.SessionToken
}

type repositoryRegistry struct {
	db *gorm.DB
}

// NewRepoRegistry ...
func NewRepoRegistry(db *gorm.DB) RepositoryRegistry {
	var r repositoryRegistry
	var once sync.Once

	once.Do(func() {
		r = repositoryRegistry{db: db}
	})

	return r
}

// func (r repositoryRegistry) User() repository.User {
// 	var once sync.Once
// 	var userRepo repository.User

// 	once.Do(func() {
// 		userRepo = repository.NewUserRepository(r.db)
// 	})

// 	return userRepo
// }

// func (r repositoryRegistry) SessionToken() repository.SessionToken {
// 	var once sync.Once
// 	var sessionTokenRepo repository.SessionToken

// 	once.Do(func() {
// 		sessionTokenRepo = repository.NewSessionToken(&http.Client{Transport: &http.Transport{}})
// 	})
// 	return sessionTokenRepo
// }
