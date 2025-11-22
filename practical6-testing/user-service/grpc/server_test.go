package grpc

import (
	"context"
	"testing"
	"user-service/database"
	"user-service/models"
	userv1 "user-service/proto/userv1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	return db
}

func teardownTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.Close()
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	server := NewUserServer()

	tests := []struct {
		name        string
		request     *userv1.CreateUserRequest
		wantErr     bool
		expectedMsg string
	}{
		{
			name: "successful user creation",
			request: &userv1.CreateUserRequest{
				Name:        "John Doe",
				Email:       "john@example.com",
				IsCafeOwner: false,
			},
			wantErr: false,
		},
		{
			name: "create cafe owner",
			request: &userv1.CreateUserRequest{
				Name:        "Jane Owner",
				Email:       "jane@cafeshop.com",
				IsCafeOwner: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := server.CreateUser(ctx, tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, resp.User.Id)
				assert.Equal(t, tt.request.Name, resp.User.Name)
				assert.Equal(t, tt.request.Email, resp.User.Email)
				assert.Equal(t, tt.request.IsCafeOwner, resp.User.IsCafeOwner)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	server := NewUserServer()
	ctx := context.Background()

	// Create a test user
	createResp, err := server.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:        "Test User",
		Email:       "test@example.com",
		IsCafeOwner: false,
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		userID      uint32
		wantErr     bool
		expectedErr codes.Code
	}{
		{
			name:    "get existing user",
			userID:  createResp.User.Id,
			wantErr: false,
		},
		{
			name:        "get non-existent user",
			userID:      9999,
			wantErr:     true,
			expectedErr: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetUser(ctx, &userv1.GetUserRequest{Id: tt.userID})

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedErr, st.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.userID, resp.User.Id)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	server := NewUserServer()
	ctx := context.Background()

	// Create multiple users
	users := []struct {
		name  string
		email string
	}{
		{"User 1", "user1@example.com"},
		{"User 2", "user2@example.com"},
		{"User 3", "user3@example.com"},
	}

	for _, u := range users {
		_, err := server.CreateUser(ctx, &userv1.CreateUserRequest{
			Name:  u.name,
			Email: u.email,
		})
		require.NoError(t, err)
	}

	// Get all users
	resp, err := server.GetUsers(ctx, &userv1.GetUsersRequest{})
	require.NoError(t, err)
	assert.Len(t, resp.Users, len(users))
}
