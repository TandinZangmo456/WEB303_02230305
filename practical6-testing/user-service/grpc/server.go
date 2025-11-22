package grpc

import (
	"context"
	"user-service/database"
	"user-service/models"
	userv1 "user-service/proto/userv1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userv1.UnimplementedUserServiceServer
}

func NewUserServer() *UserServer {
	return &UserServer{}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	user := models.User{
		Name:        req.Name,
		Email:       req.Email,
		IsCafeOwner: req.IsCafeOwner,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &userv1.CreateUserResponse{
		User: &userv1.User{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			IsCafeOwner: user.IsCafeOwner,
			CreatedAt:   user.CreatedAt.String(),
			UpdatedAt:   user.UpdatedAt.String(),
		},
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.GetUserResponse, error) {
	var user models.User
	if err := database.DB.First(&user, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &userv1.GetUserResponse{
		User: &userv1.User{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			IsCafeOwner: user.IsCafeOwner,
			CreatedAt:   user.CreatedAt.String(),
			UpdatedAt:   user.UpdatedAt.String(),
		},
	}, nil
}

func (s *UserServer) GetUsers(ctx context.Context, req *userv1.GetUsersRequest) (*userv1.GetUsersResponse, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	var pbUsers []*userv1.User
	for _, user := range users {
		pbUsers = append(pbUsers, &userv1.User{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			IsCafeOwner: user.IsCafeOwner,
			CreatedAt:   user.CreatedAt.String(),
			UpdatedAt:   user.UpdatedAt.String(),
		})
	}

	return &userv1.GetUsersResponse{Users: pbUsers}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error) {
	var user models.User
	if err := database.DB.First(&user, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	user.Name = req.Name
	user.Email = req.Email
	user.IsCafeOwner = req.IsCafeOwner

	if err := database.DB.Save(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &userv1.UpdateUserResponse{
		User: &userv1.User{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			IsCafeOwner: user.IsCafeOwner,
			CreatedAt:   user.CreatedAt.String(),
			UpdatedAt:   user.UpdatedAt.String(),
		},
	}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error) {
	result := database.DB.Delete(&models.User{}, req.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &userv1.DeleteUserResponse{Success: true}, nil
}
