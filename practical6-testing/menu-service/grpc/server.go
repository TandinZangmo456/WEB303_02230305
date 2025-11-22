package grpc

import (
	"context"
	"menu-service/database"
	"menu-service/models"
	menuv1 "menu-service/proto/menuv1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MenuServer struct {
	menuv1.UnimplementedMenuServiceServer
}

func NewMenuServer() *MenuServer {
	return &MenuServer{}
}

func (s *MenuServer) CreateMenuItem(ctx context.Context, req *menuv1.CreateMenuItemRequest) (*menuv1.CreateMenuItemResponse, error) {
	menuItem := models.MenuItem{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	}

	if err := database.DB.Create(&menuItem).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create menu item: %v", err)
	}

	return &menuv1.CreateMenuItemResponse{
		MenuItem: &menuv1.MenuItem{
			Id:          uint32(menuItem.ID),
			Name:        menuItem.Name,
			Description: menuItem.Description,
			Price:       menuItem.Price,
			CreatedAt:   menuItem.CreatedAt.String(),
			UpdatedAt:   menuItem.UpdatedAt.String(),
		},
	}, nil
}

func (s *MenuServer) GetMenuItem(ctx context.Context, req *menuv1.GetMenuItemRequest) (*menuv1.GetMenuItemResponse, error) {
	var menuItem models.MenuItem
	if err := database.DB.First(&menuItem, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "menu item not found")
	}

	return &menuv1.GetMenuItemResponse{
		MenuItem: &menuv1.MenuItem{
			Id:          uint32(menuItem.ID),
			Name:        menuItem.Name,
			Description: menuItem.Description,
			Price:       menuItem.Price,
			CreatedAt:   menuItem.CreatedAt.String(),
			UpdatedAt:   menuItem.UpdatedAt.String(),
		},
	}, nil
}

func (s *MenuServer) GetMenuItems(ctx context.Context, req *menuv1.GetMenuItemsRequest) (*menuv1.GetMenuItemsResponse, error) {
	var menuItems []models.MenuItem
	if err := database.DB.Find(&menuItems).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get menu items: %v", err)
	}

	var pbMenuItems []*menuv1.MenuItem
	for _, item := range menuItems {
		pbMenuItems = append(pbMenuItems, &menuv1.MenuItem{
			Id:          uint32(item.ID),
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			CreatedAt:   item.CreatedAt.String(),
			UpdatedAt:   item.UpdatedAt.String(),
		})
	}

	return &menuv1.GetMenuItemsResponse{MenuItems: pbMenuItems}, nil
}

func (s *MenuServer) UpdateMenuItem(ctx context.Context, req *menuv1.UpdateMenuItemRequest) (*menuv1.UpdateMenuItemResponse, error) {
	var menuItem models.MenuItem
	if err := database.DB.First(&menuItem, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "menu item not found")
	}

	menuItem.Name = req.Name
	menuItem.Description = req.Description
	menuItem.Price = req.Price

	if err := database.DB.Save(&menuItem).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update menu item: %v", err)
	}

	return &menuv1.UpdateMenuItemResponse{
		MenuItem: &menuv1.MenuItem{
			Id:          uint32(menuItem.ID),
			Name:        menuItem.Name,
			Description: menuItem.Description,
			Price:       menuItem.Price,
			CreatedAt:   menuItem.CreatedAt.String(),
			UpdatedAt:   menuItem.UpdatedAt.String(),
		},
	}, nil
}

func (s *MenuServer) DeleteMenuItem(ctx context.Context, req *menuv1.DeleteMenuItemRequest) (*menuv1.DeleteMenuItemResponse, error) {
	result := database.DB.Delete(&models.MenuItem{}, req.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete menu item: %v", result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "menu item not found")
	}

	return &menuv1.DeleteMenuItemResponse{Success: true}, nil
}
