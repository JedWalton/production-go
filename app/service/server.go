package service

import (
	"production-go/data"
)

type ServiceContainer struct {
	UserService *UserService
	// Add other services here as needed
}

func NewServiceContainer(pg *data.PostgreSQL) *ServiceContainer {
	return &ServiceContainer{
		UserService: NewUserService(pg),
		// Initialize other services here
	}
}
