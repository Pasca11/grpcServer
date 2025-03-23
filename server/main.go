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
	"google.golang.org/grpc/credentials/insecure"
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
	Status          string
	CreatedAt       time.Time
}

func NewServer() *server {
	return &server{
		orders: make(map[string]*Order),
	}
}

func (s *server) CreateOrder(ctx context.Context, req *pb.OrderRequest) (*pb.OrderResponse, error) {
	orderId := uuid.New().String()
	return &pb.OrderResponse{
		OrderId:               orderId,
		EstimatedDeliveryTime: time.Now().Add(2 * time.Hour).Format(time.RFC3339),
	}, nil
}

func (s *server) GetOrderStatus(ctx context.Context, req *pb.OrderStatusRequest) (*pb.OrderStatusResponse, error) {
	_, exists := s.orders[req.OrderId]
	if !exists {
		return nil, status.Error(codes.NotFound, "заказ не найден")
	}

	return &pb.OrderStatusResponse{
		Status:            "pending",
		StatusDescription: "Заказ обрабатывается",
	}, nil
}

func (s *server) GetAllOrders(ctx context.Context, req *pb.GetAllOrdersRequest) (*pb.GetAllOrdersResponse, error) {
	orders := []*pb.Order{
		{
			OrderId:           "1",
			CustomerName:      "John Doe",
			DeliveryAddress:   "123 Main St",
			BottlesCount:      10,
			PhoneNumber:       "1234567890",
			Status:            "pending",
			StatusDescription: "Заказ обрабатывается",
		},
		{
			OrderId:           "2",
			CustomerName:      "Jane Smith",
			DeliveryAddress:   "456 Oak Ave",
			BottlesCount:      5,
			PhoneNumber:       "0987654321",
			Status:            "confirmed",
			StatusDescription: "Заказ подтвержден",
		},
		{
			OrderId:           "3",
			CustomerName:      "Alice Johnson",
			DeliveryAddress:   "789 Pine Rd",
			BottlesCount:      20,
			PhoneNumber:       "1122334455",
			Status:            "in_delivery",
			StatusDescription: "Заказ в пути",
		},
		{
			OrderId:           "4",
			CustomerName:      "Bob Brown",
			DeliveryAddress:   "101 Maple St",
			BottlesCount:      15,
			PhoneNumber:       "9988776655",
			Status:            "delivered",
			StatusDescription: "Заказ доставлен",
		},
		{
			OrderId:           "5",
			CustomerName:      "Charlie Davis",
			DeliveryAddress:   "222 Cedar Ave",
			BottlesCount:      30,
			PhoneNumber:       "5554443322",
			Status:            "cancelled",
			StatusDescription: "Заказ отменен",
		},
	}

	return &pb.GetAllOrdersResponse{Orders: orders}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	pb.RegisterWaterDeliveryServiceServer(s, NewServer())
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
