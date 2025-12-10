package handler

import (
	"social-network/services/users/internal/application"
	pb "social-network/shared/gen-go/users"
)

// Holds Client conns, services and handler funcs
type UsersHandler struct {
	pb.UserServiceServer
	Application *application.Application
}

func NewUsersHanlder(service *application.Application) *UsersHandler {
	return &UsersHandler{
		Application: service,
	}
}
