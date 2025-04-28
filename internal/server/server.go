package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/handmade-jewellery/user-service/internal/app/user"
	pb "github.com/handmade-jewellery/user-service/pkg/api/user-service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
	"net/http"
)

type Server struct {
	userServiceServer *user.Service
	grpcPort          string
	httpPort          string
}

func NewServer(userServiceServer *user.Service) *Server {
	return &Server{
		userServiceServer: userServiceServer,
	}
}

func (s *Server) Run() error {
	//1. Запуск gRPC сервера порт 8081
	err := s.runGRPC()
	if err != nil {
		log.Printf("failed to run gRPC server: %v", err)
		return err
	}

	// 2. Запуск http сервера порт 8080
	err = s.runHTTP()
	if err != nil {
		log.Printf("failed to run http server: %v", err)
		return err
	}

	return nil
}

func (s *Server) runGRPC() error {
	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, s.userServiceServer)

	lis, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Printf("Failed to listen: %v", err)
		return err
	}

	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Printf("gRPC server failed: %v", err)
			s.stopGRPC(grpcServer)
		}
	}()

	log.Println("gRPC server started on :8080")
	return nil
}

func (s *Server) stopGRPC(grpcServer *grpc.Server) {
	// graceful shutdown
	grpcServer.GracefulStop()
	log.Println("gRPC server gracefully stopped")
}

func (s *Server) runHTTP() error {
	mux := runtime.NewServeMux()

	conn, err := grpc.NewClient("localhost:8084", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to create gRPC client: %v", err)
		return err
	}

	defer conn.Close() // Обязательно закрыть соединение после завершения

	err = pb.RegisterUserServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Printf("Failed to register user service handler: %v", err)
		return err
	}

	// 3. Запуск HTTP сервера с маршрутизацией через Chi
	router := chi.NewRouter()
	router.Handle("/", mux) // Гейтвей будет обрабатывать HTTP-запросы

	err = http.ListenAndServe(":8085", router)
	if err != nil {
		log.Printf("Failed to start HTTP server: %v", err)
		return err
	}

	log.Println("Starting HTTP server on :8080...")
	return nil
}
