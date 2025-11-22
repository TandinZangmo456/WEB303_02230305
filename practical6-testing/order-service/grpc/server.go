package grpc

import (
	"context"
	"order-service/database"
	"order-service/models"
	orderv1 "order-service/proto/orderv1"

	menuv1 "order-service/proto/menuv1"
	userv1 "order-service/proto/userv1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderServer struct {
	orderv1.UnimplementedOrderServiceServer
	UserClient userv1.UserServiceClient
	MenuClient menuv1.MenuServiceClient
}

func NewOrderServer(userClient userv1.UserServiceClient, menuClient menuv1.MenuServiceClient) *OrderServer {
	return &OrderServer{
		UserClient: userClient,
		MenuClient: menuClient,
	}
}

func (s *OrderServer) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	// Validate user exists
	_, err := s.UserClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.UserId})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "user not found: %v", err)
	}

	// Create order
	order := models.Order{
		UserID: uint(req.UserId),
		Status: "pending",
	}

	// Process order items
	for _, item := range req.Items {
		// Validate menu item exists and get price
		menuResp, err := s.MenuClient.GetMenuItem(ctx, &menuv1.GetMenuItemRequest{Id: item.MenuItemId})
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "menu item %d not found: %v", item.MenuItemId, err)
		}

		orderItem := models.OrderItem{
			MenuItemID:   uint(item.MenuItemId),
			MenuItemName: menuResp.MenuItem.Name,
			Quantity:     uint(item.Quantity),
			Price:        menuResp.MenuItem.Price, // Snapshot the price
		}
		order.OrderItems = append(order.OrderItems, orderItem)
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create order: %v", err)
	}

	// Build response
	var pbOrderItems []*orderv1.OrderItem
	for _, item := range order.OrderItems {
		pbOrderItems = append(pbOrderItems, &orderv1.OrderItem{
			Id:           uint32(item.ID),
			MenuItemId:   uint32(item.MenuItemID),
			MenuItemName: item.MenuItemName,
			Quantity:     uint32(item.Quantity),
			Price:        item.Price,
		})
	}

	return &orderv1.CreateOrderResponse{
		Order: &orderv1.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			Status:     order.Status,
			OrderItems: pbOrderItems,
			CreatedAt:  order.CreatedAt.String(),
			UpdatedAt:  order.UpdatedAt.String(),
		},
	}, nil
}

func (s *OrderServer) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	var order models.Order
	if err := database.DB.Preload("OrderItems").First(&order, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	var pbOrderItems []*orderv1.OrderItem
	for _, item := range order.OrderItems {
		pbOrderItems = append(pbOrderItems, &orderv1.OrderItem{
			Id:           uint32(item.ID),
			MenuItemId:   uint32(item.MenuItemID),
			MenuItemName: item.MenuItemName,
			Quantity:     uint32(item.Quantity),
			Price:        item.Price,
		})
	}

	return &orderv1.GetOrderResponse{
		Order: &orderv1.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			Status:     order.Status,
			OrderItems: pbOrderItems,
			CreatedAt:  order.CreatedAt.String(),
			UpdatedAt:  order.UpdatedAt.String(),
		},
	}, nil
}

func (s *OrderServer) GetOrders(ctx context.Context, req *orderv1.GetOrdersRequest) (*orderv1.GetOrdersResponse, error) {
	var orders []models.Order
	if err := database.DB.Preload("OrderItems").Find(&orders).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
	}

	var pbOrders []*orderv1.Order
	for _, order := range orders {
		var pbOrderItems []*orderv1.OrderItem
		for _, item := range order.OrderItems {
			pbOrderItems = append(pbOrderItems, &orderv1.OrderItem{
				Id:           uint32(item.ID),
				MenuItemId:   uint32(item.MenuItemID),
				MenuItemName: item.MenuItemName,
				Quantity:     uint32(item.Quantity),
				Price:        item.Price,
			})
		}

		pbOrders = append(pbOrders, &orderv1.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			Status:     order.Status,
			OrderItems: pbOrderItems,
			CreatedAt:  order.CreatedAt.String(),
			UpdatedAt:  order.UpdatedAt.String(),
		})
	}

	return &orderv1.GetOrdersResponse{Orders: pbOrders}, nil
}

func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *orderv1.UpdateOrderStatusRequest) (*orderv1.UpdateOrderStatusResponse, error) {
	var order models.Order
	if err := database.DB.Preload("OrderItems").First(&order, req.Id).Error; err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found")
	}

	order.Status = req.Status
	if err := database.DB.Save(&order).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order status: %v", err)
	}

	var pbOrderItems []*orderv1.OrderItem
	for _, item := range order.OrderItems {
		pbOrderItems = append(pbOrderItems, &orderv1.OrderItem{
			Id:           uint32(item.ID),
			MenuItemId:   uint32(item.MenuItemID),
			MenuItemName: item.MenuItemName,
			Quantity:     uint32(item.Quantity),
			Price:        item.Price,
		})
	}

	return &orderv1.UpdateOrderStatusResponse{
		Order: &orderv1.Order{
			Id:         uint32(order.ID),
			UserId:     uint32(order.UserID),
			Status:     order.Status,
			OrderItems: pbOrderItems,
			CreatedAt:  order.CreatedAt.String(),
			UpdatedAt:  order.UpdatedAt.String(),
		},
	}, nil
}
