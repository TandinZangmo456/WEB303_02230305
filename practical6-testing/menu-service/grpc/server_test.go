package grpc

import (
	"context"
	"testing"
	"menu-service/database"
	"menu-service/models"
	menuv1 "menu-service/proto/menuv1"

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

	err = db.AutoMigrate(&models.MenuItem{})
	require.NoError(t, err)

	return db
}

func teardownTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	sqlDB.Close()
}

func TestCreateMenuItem(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	server := NewMenuServer()

	tests := []struct {
		name    string
		request *menuv1.CreateMenuItemRequest
		wantErr bool
	}{
		{
			name: "successful menu item creation",
			request: &menuv1.CreateMenuItemRequest{
				Name:        "Cappuccino",
				Description: "Espresso with steamed milk",
				Price:       4.50,
			},
			wantErr: false,
		},
		{
			name: "create item with zero price",
			request: &menuv1.CreateMenuItemRequest{
				Name:  "Water",
				Price: 0.0,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := server.CreateMenuItem(ctx, tt.request)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotZero(t, resp.MenuItem.Id)
				assert.Equal(t, tt.request.Name, resp.MenuItem.Name)
				assert.InDelta(t, tt.request.Price, resp.MenuItem.Price, 0.001)
			}
		})
	}
}

func TestGetMenuItem(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	server := NewMenuServer()
	ctx := context.Background()

	// Create a test menu item
	createResp, err := server.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
		Name:        "Test Coffee",
		Description: "Test Description",
		Price:       3.50,
	})
	require.NoError(t, err)

	tests := []struct {
		name        string
		itemID      uint32
		wantErr     bool
		expectedErr codes.Code
	}{
		{
			name:    "get existing menu item",
			itemID:  createResp.MenuItem.Id,
			wantErr: false,
		},
		{
			name:        "get non-existent menu item",
			itemID:      9999,
			wantErr:     true,
			expectedErr: codes.NotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := server.GetMenuItem(ctx, &menuv1.GetMenuItemRequest{Id: tt.itemID})

			if tt.wantErr {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				assert.Equal(t, tt.expectedErr, st.Code())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.itemID, resp.MenuItem.Id)
			}
		})
	}
}

func TestGetMenuItems(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	server := NewMenuServer()
	ctx := context.Background()

	// Create multiple menu items
	items := []struct {
		name  string
		price float64
	}{
		{"Coffee", 2.50},
		{"Tea", 2.00},
		{"Sandwich", 5.00},
	}

	for _, item := range items {
		_, err := server.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
			Name:  item.name,
			Price: item.price,
		})
		require.NoError(t, err)
	}

	// Get all menu items
	resp, err := server.GetMenuItems(ctx, &menuv1.GetMenuItemsRequest{})
	require.NoError(t, err)
	assert.Len(t, resp.MenuItems, len(items))
}

func TestPriceHandling(t *testing.T) {
	db := setupTestDB(t)
	defer teardownTestDB(t, db)
	database.DB = db

	server := NewMenuServer()
	ctx := context.Background()

	testCases := []struct {
		name  string
		price float64
	}{
		{"integer price", 5.0},
		{"two decimal places", 5.99},
		{"very small price", 0.01},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := server.CreateMenuItem(ctx, &menuv1.CreateMenuItemRequest{
				Name:  "Test",
				Price: tc.price,
			})

			require.NoError(t, err)
			assert.InDelta(t, tc.price, resp.MenuItem.Price, 0.001)
		})
	}
}
