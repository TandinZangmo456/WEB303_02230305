package grpc

import (
	"context"
	"testing"
	"order-service/database"
	"order-service/models"
	orderv1 "order-service/proto/orderv1"
	userv1 "order-service/proto/userv1"
	menuv1 "order-service/proto/menuv1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// MockUserServiceClient simulates the user service
type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) CreateUser(ctx context.Context, req *userv1.CreateUserRequest, opts ...grpc.CallOption) (*userv1.CreateUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userv1.CreateUserResponse), args.Error(1)
}

func (m *MockUserServiceClient) GetUser(ctx context.Context, req *userv1.GetUserRequest, opts ...grpc.CallOption) (*userv1.GetUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userv1.GetUserResponse), args.Error(1)
}

func (m *MockUserServiceClient) GetUsers(ctx context.Context, req *userv1.GetUsersRequest, opts ...grpc.CallOption) (*userv1.GetUsersResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userv1.GetUsersResponse), args.Error(1)
}

func (m *MockUserServiceClient) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest, opts ...grpc.CallOption) (*userv1.UpdateUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userv1.UpdateUserResponse), args.Error(1)
}

func (m *MockUserServiceClient) DeleteUser(ctx context.Context, req *userv1.DeleteUserRequest, opts ...grpc.CallOption) (*userv1.DeleteUserResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userv1.DeleteUserResponse), args.Error(1)
}

// MockMenuServiceClient simulates the menu service
type MockMenuServiceClient struct {
	mock.Mock
}

func (m *MockMenuServiceClient) CreateMenuItem(ctx context.Context, req *menuv1.CreateMenuItemRequest, opts ...grpc.CallOption) (*menuv1.CreateMenuItemResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*menuv1.CreateMenuItemResponse), args.Error(1)
}

func (m *MockMenuServiceClient) GetMenuItem(ctx context.Context, req *menuv1.GetMenuItemRequest, opts ...grpc.CallOption) (*menuv1.GetMenuItemResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*menuv1.GetMenuItemResponse), args.Error(1)
}

func (m *MockMenuServiceClient) GetMenuItems(ctx context.Context, req *menuv1.GetMenuItemsRequest, opts ...grpc.CallOption) (*menuv1.GetMenuItemsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*menuv1.GetMenuItemsResponse), args.Error(1)
}

func (m *MockMenuServiceClient) UpdateMenuItem(ctx context.Context, req *menuv1.UpdateMenuItemRequest, opts ...grpc.CallOption) (*menuv1.UpdateMenuItemResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*menuv1.UpdateMenuItemResponse), args.Error(1)
}

func (m *MockMenuServiceClient) DeleteMenuItem(ctx context.Context, req *menuv1.DeleteMenuItemRequest, opts ...grpc.CallOption) (*menuv1.DeleteMenuItemResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*menuv1.DeleteMenuItemResponse), args.Error(1)
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Order{}, &models.OrderItem{})
	require.NoError(t, err)

	return db
}

func teardownTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.Close()
}

func TestCreateOrder_Success(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	mockUserClient := new(MockUserServiceClient)
	mockMenuClient := new(MockMenuServiceClient)

	server := &OrderServer{
		UserClient: mockUserClient,
		MenuClient: mockMenuClient,
	}

	mockUserClient.On("GetUser", mock.Anything, &userv1.GetUserRequest{Id: 1}).
		Return(&userv1.GetUserResponse{
			User: &userv1.User{Id: 1, Name: "Test User"},
		}, nil)

	mockMenuClient.On("GetMenuItem", mock.Anything, &menuv1.GetMenuItemRequest{Id: 1}).
		Return(&menuv1.GetMenuItemResponse{
			MenuItem: &menuv1.MenuItem{Id: 1, Name: "Coffee", Price: 2.50},
		}, nil)

	ctx := context.Background()
	resp, err := server.CreateOrder(ctx, &orderv1.CreateOrderRequest{
		UserId: 1,
		Items: []*orderv1.OrderItemRequest{
			{MenuItemId: 1, Quantity: 2},
		},
	})

	require.NoError(t, err)
	assert.NotZero(t, resp.Order.Id)
	assert.Equal(t, uint32(1), resp.Order.UserId)
	assert.Len(t, resp.Order.OrderItems, 1)
	assert.InDelta(t, 2.50, resp.Order.OrderItems[0].Price, 0.001)

	mockUserClient.AssertExpectations(t)
	mockMenuClient.AssertExpectations(t)
}

func TestCreateOrder_InvalidUser(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	mockUserClient := new(MockUserServiceClient)
	mockMenuClient := new(MockMenuServiceClient)

	server := &OrderServer{
		UserClient: mockUserClient,
		MenuClient: mockMenuClient,
	}

	mockUserClient.On("GetUser", mock.Anything, &userv1.GetUserRequest{Id: 999}).
		Return(nil, status.Errorf(codes.NotFound, "user not found"))

	ctx := context.Background()
	resp, err := server.CreateOrder(ctx, &orderv1.CreateOrderRequest{
		UserId: 999,
		Items: []*orderv1.OrderItemRequest{
			{MenuItemId: 1, Quantity: 1},
		},
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "user not found")

	mockUserClient.AssertExpectations(t)
}

func TestCreateOrder_InvalidMenuItem(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	mockUserClient := new(MockUserServiceClient)
	mockMenuClient := new(MockMenuServiceClient)

	server := &OrderServer{
		UserClient: mockUserClient,
		MenuClient: mockMenuClient,
	}

	mockUserClient.On("GetUser", mock.Anything, &userv1.GetUserRequest{Id: 1}).
		Return(&userv1.GetUserResponse{
			User: &userv1.User{Id: 1, Name: "Test User"},
		}, nil)

	mockMenuClient.On("GetMenuItem", mock.Anything, &menuv1.GetMenuItemRequest{Id: 999}).
		Return(nil, status.Errorf(codes.NotFound, "menu item not found"))

	ctx := context.Background()
	resp, err := server.CreateOrder(ctx, &orderv1.CreateOrderRequest{
		UserId: 1,
		Items: []*orderv1.OrderItemRequest{
			{MenuItemId: 999, Quantity: 1},
		},
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	require.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "menu item 999 not found")

	mockUserClient.AssertExpectations(t)
	mockMenuClient.AssertExpectations(t)
}
