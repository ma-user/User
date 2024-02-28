package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var addr string = "0.0.0.0:50051"

type userServiceServer struct {
	DB *gorm.DB
	pb.UserServiceServer
}

type User struct {
	gorm.Model
	id         int64
	First_name string
	Last_name  string
	Age        int32
	Token      string
}

func initialize(dsn string) *gorm.DB {
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to db", err)
	}

	DB.AutoMigrate(&User{})

	fmt.Println("Connected to DB successfully!")
	return DB
}

func (s *userServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := req.GetUser()

	if user == nil || user.GetFirstName() == "" || user.GetLastName() == "" || user.GetAge() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user data")
	}

	token := uuid.New().String()

	users := User{
		id:         user.Id,
		First_name: user.FirstName,
		Last_name:  user.LastName,
		Age:        user.Age,
		Token:      token,
	}

	if s.DB == nil {
		return nil, status.Error(codes.Internal, "Database connection is nil")
	}

	result := s.DB.Create(&users)

	if result.Error != nil {
		return nil, status.Error(codes.Internal, result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return nil, errors.New("cannot create user successfully")
	}

	response := &pb.CreateUserResponse{
		User:    user,
		Token:   users.Token,
		Message: "Created user successfully",
	}

	return response, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	id := req.Id
	token := req.Token

	var user *pb.User

	s.DB.First(&user, id)

	if user.Id == 0 {
		return nil, errors.New("user not found")
	}

	if user.Token != token {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	response := &pb.GetUserResponse{
		User: user,
	}

	return response, nil
}

func (s *userServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	usr := req.User
	id := req.Id
	token := req.Token

	if usr == nil || usr.GetFirstName() == "" || usr.GetLastName() == "" || usr.GetAge() <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid user data")
	}

	var user *pb.User

	s.DB.First(&user, id)

	if user.Id == 0 {
		return nil, errors.New("user not found")
	}

	if user.Token != token {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	user.FirstName = usr.FirstName
	user.LastName = usr.LastName
	user.Age = usr.Age

	s.DB.Save(&user)

	response := &pb.UpdateUserResponse{
		User:    user,
		Message: "User successfully updated",
	}

	return response, nil
}

func main() {
	dataSourceName := "user=postgres password=pgpswd dbname=UserDB host=localhost port=5433 sslmode=disable"

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", err)
	}

	log.Printf("listening on %s\n", addr)

	grpcServer := grpc.NewServer()
	db := initialize(dataSourceName)

	pb.RegisterUserServiceServer(grpcServer, &userServiceServer{DB: db})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
