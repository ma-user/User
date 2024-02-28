package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
)

var (
	grpcServerAddr = "localhost:50051"
	httpServerAddr = ":8080"
)

type UserDetails struct {
	Id         int64  `json: "id"`
	First_name string `json: "first_name"`
	Last_name  string `json: "last_name"`
	Age        int32  `json: "age"`
}

func Create(client pb.UserServiceClient, w http.ResponseWriter, r *http.Request) {

	var usr UserDetails

	err := json.NewDecoder(r.Body).Decode(&usr)

	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	user := &pb.User{
		Id:        usr.Id,
		FirstName: usr.First_name,
		LastName:  usr.Last_name,
		Age:       usr.Age,
	}

	res, err := client.CreateUser(context.Background(), &pb.CreateUserRequest{
		User: user,
	})

	if err != nil {
		log.Fatal("error while creating user %v", err)
	}

	json.NewEncoder(w).Encode(struct {
		User    *pb.User `json:"user,omitempty"`
		Token   string   `json:"token"`
		Message string   ` json:"message"`
	}{
		User:    res.User,
		Token:   res.Token,
		Message: res.Message,
	})
}

func GetUser(client pb.UserServiceClient, w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)

	userId, err := strconv.Atoi(param["id"])
	if err != nil {
		panic(err)
	}

	bearerToken := extractBearerToken(r)

	if bearerToken == "" {
		http.Error(w, "Unauthorized: Bearer token not provided", http.StatusUnauthorized)
		return
	}

	res, err := client.GetUser(context.Background(), &pb.GetUserRequest{
		Id:    int64(userId),
		Token: bearerToken,
	})

	if err != nil {
		statusErr, ok := status.FromError(err)

		if ok {
			switch statusErr.Code() {
			case codes.Unauthenticated:
				http.Error(w, "Unauthorized: Invalid Bearer token", http.StatusUnauthorized)
				return
			default:
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	usr := UserDetails{
		Id:         res.User.Id,
		First_name: res.User.FirstName,
		Last_name:  res.User.LastName,
		Age:        res.User.Age,
	}

	json.NewEncoder(w).Encode(usr)

}

func UpdateUser(client pb.UserServiceClient, w http.ResponseWriter, r *http.Request) {
	var usr UserDetails

	param := mux.Vars(r)

	userId, err := strconv.Atoi(param["id"])

	if err != nil {
		panic(err)
	}

	bearerToken := extractBearerToken(r)

	if bearerToken == "" {
		http.Error(w, "Unauthorized: Bearer token not provided", http.StatusUnauthorized)
		return
	}

	json.NewDecoder(r.Body).Decode(&usr)

	user := &pb.User{
		Id:        usr.Id,
		FirstName: usr.First_name,
		LastName:  usr.Last_name,
		Age:       usr.Age,
	}

	res, err := client.UpdateUser(context.Background(), &pb.UpdateUserRequest{
		User:  user,
		Id:    int64(userId),
		Token: bearerToken,
	})

	if err != nil {
		statusErr, ok := status.FromError(err)

		if ok {
			switch statusErr.Code() {
			case codes.Unauthenticated:
				http.Error(w, "Unauthorized: Invalid Bearer token", http.StatusUnauthorized)
				return
			default:
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		User    *pb.User `json:"user,omitempty"`
		Message string   ` json:"message"`
	}{
		User:    res.User,
		Message: res.Message,
	})
}

func extractBearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		authParts := strings.Split(authHeader, " ")
		if len(authParts) == 2 && strings.ToLower(authParts[0]) == "bearer" {
			return authParts[1]
		}
	}
	return ""
}

func main() {
	conn, err := grpc.Dial(grpcServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	router := mux.NewRouter()

	router.HandleFunc("/user", func(writer http.ResponseWriter, req *http.Request) {
		Create(client, writer, req)
	}).Methods("POST")

	router.HandleFunc("/user/{id}", func(writer http.ResponseWriter, req *http.Request) {
		GetUser(client, writer, req)
	}).Methods("GET")

	router.HandleFunc("/user/{id}", func(writer http.ResponseWriter, req *http.Request) {
		UpdateUser(client, writer, req)
	}).Methods("PUT")

	log.Println("HTTP Server listening on", httpServerAddr)
	http.ListenAndServe(httpServerAddr, router)
}
