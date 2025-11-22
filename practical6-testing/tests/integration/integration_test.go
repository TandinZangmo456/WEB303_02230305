package integration

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"
	"time"

	usergrpc "user-service/grpc"
	usermodels "user-service/models"
	userdatabase "user-service/database"
	userv1 "user-service/proto/userv1"

	menugrpc "menu-service/grpc"
	menumodels "menu-service/models"
	menudatabase "menu-service/database"
	menuv1 "menu-service/proto/menuv1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

func bufDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}

func setupUserService(t *testing.T) (*bufconn.Listener, *grpc.Server) {
	// Use unique database name for each test
	dbName := fmt.Sprintf("file:test_%d.db?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&usermodels.User{})
	require.NoError(t, err)

	userdatabase.DB = db

	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	userv1.RegisterUserServiceServer(server, usergrpc.NewUserServer())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Printf("User server stopped: %v", err)
		}
	}()
	
	// Give server time to start
	time.Sleep(50 * time.Millisecond)
	
	return listener, server
}

func setupMenuService(t *testing.T) (*bufconn.Listener, *grpc.Server) {
	// Use unique database name for each test
	dbName := fmt.Sprintf("file:test_%d.db?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&menumodels.MenuItem{})
	require.NoError(t, err)

	menudatabase.DB = db

	listener := bufconn.Listen(bufSize)
	server := grpc.NewServer()
	menuv1.RegisterMenuServiceServer(server, menugrpc.NewMenuServer())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Printf("Menu server stopped: %v", err)
		}
	}()
	
	// Give server time to start
	time.Sleep(50 * time.Millisecond)
	
	return listener, server
}

func TestIntegration_CreateUser(t *testing.T) {
	userListener, userServer := setupUserService(t)
	defer func() {
		userServer.Stop()
		userListener.Close()
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(userListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := userv1.NewUserServiceClient(conn)

	resp, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:        "Integration User",
		Email:       "integration@test.com",
		IsCafeOwner: false,
	})

	require.NoError(t, err)
	assert.NotZero(t, resp.User.Id)
	assert.Equal(t, "Integration User", resp.User.Name)
	assert.Equal(t, "integration@test.com", resp.User.Email)
}

func TestIntegration_CreateMenuItem(t *testing.T) {
	menuListener, menuServer := setupMenuService(t)
	defer func() {
		menuServer.Stop()
		menuListener.Close()
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(menuListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := menuv1.NewMenuServiceClient(conn)

	resp, err := client.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
		Name:        "Coffee",
		Description: "Hot coffee",
		Price:       2.50,
	})

	require.NoError(t, err)
	assert.NotZero(t, resp.MenuItem.Id)
	assert.Equal(t, "Coffee", resp.MenuItem.Name)
	assert.InDelta(t, 2.50, resp.MenuItem.Price, 0.001)
}

func TestIntegration_UserAndMenuServices(t *testing.T) {
	userListener, userServer := setupUserService(t)
	defer func() {
		userServer.Stop()
		userListener.Close()
	}()

	menuListener, menuServer := setupMenuService(t)
	defer func() {
		menuServer.Stop()
		menuListener.Close()
	}()

	ctx := context.Background()

	// Connect to user service
	userConn, err := grpc.DialContext(ctx, "user-bufnet",
		grpc.WithContextDialer(bufDialer(userListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer userConn.Close()

	userClient := userv1.NewUserServiceClient(userConn)

	// Connect to menu service
	menuConn, err := grpc.DialContext(ctx, "menu-bufnet",
		grpc.WithContextDialer(bufDialer(menuListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer menuConn.Close()

	menuClient := menuv1.NewMenuServiceClient(menuConn)

	// Create a user
	userResp, err := userClient.CreateUser(ctx, &userv1.CreateUserRequest{
		Name:        "Test User",
		Email:       "test@example.com",
		IsCafeOwner: false,
	})
	require.NoError(t, err)
	assert.NotZero(t, userResp.User.Id)

	// Create menu items
	item1, err := menuClient.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
		Name:        "Coffee",
		Description: "Hot coffee",
		Price:       2.50,
	})
	require.NoError(t, err)

	item2, err := menuClient.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
		Name:        "Sandwich",
		Description: "Ham sandwich",
		Price:       5.00,
	})
	require.NoError(t, err)

	// Verify both services are working
	assert.NotZero(t, userResp.User.Id)
	assert.NotZero(t, item1.MenuItem.Id)
	assert.NotZero(t, item2.MenuItem.Id)
	
	// Get all users
	usersResp, err := userClient.GetUsers(ctx, &userv1.GetUsersRequest{})
	require.NoError(t, err)
	assert.Len(t, usersResp.Users, 1)
	
	// Get all menu items
	itemsResp, err := menuClient.GetMenuItems(ctx, &menuv1.GetMenuItemsRequest{})
	require.NoError(t, err)
	assert.Len(t, itemsResp.MenuItems, 2)
}

func TestIntegration_MultipleUsers(t *testing.T) {
	userListener, userServer := setupUserService(t)
	defer func() {
		userServer.Stop()
		userListener.Close()
	}()

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet",
		grpc.WithContextDialer(bufDialer(userListener)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := userv1.NewUserServiceClient(conn)

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
		resp, err := client.CreateUser(ctx, &userv1.CreateUserRequest{
			Name:  u.name,
			Email: u.email,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.User.Id)
	}

	// Get all users
	usersResp, err := client.GetUsers(ctx, &userv1.GetUsersRequest{})
	require.NoError(t, err)
	assert.Len(t, usersResp.Users, 3)
}
