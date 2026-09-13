package handler

import (
	"context"

	"github.com/amodemoli/microservices/exercise/user/internal/helpers/codes/response"
	"github.com/amodemoli/microservices/exercise/user/internal/helpers/codes/status"
	"github.com/amodemoli/microservices/exercise/user/protobuf"
)

// make handler for handle protobuf request codes in go
type Handler struct {
	protobuf.UnimplementedUserServer
}

// GetUser method, handle's GetUser request on protobuf
func (h *Handler) GetUser(ctx context.Context, req *protobuf.GetUserRequest) (*protobuf.UserResponse, error) {
	// fake user-id check
	if req.Id != 1 {
		return &protobuf.UserResponse{Response: &protobuf.ResponseModel{
			Error:  "user not-found",
			Status: string(status.NotFound),
			Code:   string(response.UserNotFound),
		}}, nil
	}
	// and return fake user information
	return &protobuf.UserResponse{
		Response: &protobuf.ResponseModel{
			Status: string(status.Ok),
			Code:   string(response.Ok),
		},
		User: &protobuf.UserModel{
			Id:          req.Id,
			Email:       "me@demolition.ir",
			Username:    "demolition",
			Displayname: "amir",
		},
	}, nil
}

// CreateUser method, handle's CreateUser request on protobuf
func (h *Handler) CreateUser(ctx context.Context, req *protobuf.CreateUserRequest) (*protobuf.UserResponse, error) {

	// return fake response
	return &protobuf.UserResponse{
		// need to add status and ...
		User: &protobuf.UserModel{
			Id:          2,
			Email:       req.Email,
			Username:    req.Username,
			Displayname: req.Displayname,
		},
	}, nil
}
