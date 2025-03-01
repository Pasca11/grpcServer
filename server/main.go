package main

import (
	"context"
	"log"
	"net"
	"time"

	pb "github.com/Pasca11/grpcServer/proto/gen"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedWaterDeliveryServiceServer
	orders map[string]*Order
}

type Order struct {
	CustomerName    string
	DeliveryAddress string
	BottlesCount    int32
	PhoneNumber     string
	Status          pb.OrderStatusResponse_Status
	CreatedAt       time.Time
}

func NewServer() *server {
	return &server{
		orders: make(map[string]*Order),
	}
}

func (s *server) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	if req.BottlesCount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "количество бутылок должно быть больше 0")
	}

	orderId := uuid.New().String()
	s.orders[orderId] = &Order{
		CustomerName:    req.CustomerName,
		DeliveryAddress: req.DeliveryAddress,
		BottlesCount:    req.BottlesCount,
		PhoneNumber:     req.PhoneNumber,
		Status:          pb.OrderStatusResponse_PENDING,
		CreatedAt:       time.Now(),
	}

	return &pb.OrderResponse{
		OrderId:               orderId,
		EstimatedDeliveryTime: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, req *pb.OrderStatusRequest) (*pb.OrderStatusResponse, error) {
	order, exists := s.orders[req.OrderId]
	if !exists {
		return nil, status.Error(codes.NotFound, "заказ не найден")
	}

	return &pb.OrderStatusResponse{
		Status:            order.Status,
		StatusDescription: getStatusDescription(order.Status),
	}, nil
}

func getStatusDescription(status pb.OrderStatusResponse_Status) string {
	switch status {
	case pb.OrderStatusResponse_PENDING:
		return "Заказ обрабатывается"
	case pb.OrderStatusResponse_CONFIRMED:
		return "Заказ подтвержден"
	case pb.OrderStatusResponse_IN_DELIVERY:
		return "Заказ в пути"
	case pb.OrderStatusResponse_DELIVERED:
		return "Заказ доставлен"
	case pb.OrderStatusResponse_CANCELLED:
		return "Заказ отменен"
	default:
		return "Неизвестный статус"
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterWaterDeliveryServiceServer(s, NewServer())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
