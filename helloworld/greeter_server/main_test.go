package main

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/status"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateUser_InvalidUserData_UserNil_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.CreateUserRequest{
		User: nil,
	}

	resp, err := server.CreateUser(context.Background(), req)

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")

	statusErr, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")

	expectedErrorMsg := "Invalid user data"
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestCreateUser_InvalidUserData_FirstNameEmpty_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.CreateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "",
			LastName:  "Name",
			Age:       10,
		},
	}

	resp, err := server.CreateUser(context.Background(), req)

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")

	statusErr, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")

	expectedErrorMsg := "Invalid user data"
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestCreateUser_InvalidUserData_LastNameEmpty_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.CreateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "User",
			LastName:  "",
			Age:       10,
		},
	}

	resp, err := server.CreateUser(context.Background(), req)

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")

	statusErr, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")

	expectedErrorMsg := "Invalid user data"
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestCreateUser_InvalidUserData_NegativeAge_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.CreateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "User",
			LastName:  "Name",
			Age:       -10,
		},
	}

	resp, err := server.CreateUser(context.Background(), req)

	statusErr, ok := status.FromError(err)

	expectedErrorMsg := "Invalid user data"

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestCreateUser_InvalidUserData_ZeroAge_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.CreateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "User",
			LastName:  "Name",
			Age:       0,
		},
	}

	resp, err := server.CreateUser(context.Background(), req)

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")

	statusErr, ok := status.FromError(err)

	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")

	expectedErrorMsg := "Invalid user data"
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestCreateUser_success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"}).AddRow("1", "Cool", "Kid", 10)
	mock.ExpectQuery("INSERT").WillReturnRows(rows)

	mock.ExpectCommit()

	req := &pb.CreateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "Cool",
			LastName:  "Kid",
			Age:       10,
		},
	}

	resp, err := server.CreateUser(context.Background(), req)

	assert.NoError(t, err, "Unexpected error in CreatUser")
	assert.NotNil(t, resp, "Expected non-nil response")
	assert.Equal(t, "Created user successfully", resp.Message, "Unexpected response message")
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, "Cool", resp.User.FirstName)
	assert.Equal(t, "Kid", resp.User.LastName)
	assert.Equal(t, int32(10), resp.User.Age)
	assert.NotNil(t, resp.User.Token)
}

func TestCreateUser_ZeroRowsAffected_throwsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age"})
	mock.ExpectQuery("INSERT").WillReturnRows(rows)

	mock.ExpectCommit()

	req := &pb.CreateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "User",
			LastName:  "Name",
			Age:       10,
		},
	}

	resp, err := server.CreateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "cannot create user successfully")
}

func TestGetUser_GotUserIDZeroFromDB_UserNotFound_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "token"}).AddRow(0, "Cool", "Kid", 12, "valid_token")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	req := &pb.GetUserRequest{
		Id:    1,
		Token: "valid_token",
	}

	resp, err := server.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.EqualError(t, err, "user not found")
}

func TestGetUser_InvalidToken_NotAuthenticated_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "token"}).AddRow(1, "Cool", "Kid", 12, "valid_token")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	req := &pb.GetUserRequest{
		Id:    1,
		Token: "invalid_token",
	}

	resp, err := server.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestGetUser_success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "token"}).AddRow(1, "Cool", "Kid", 12, "valid_token")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	req := &pb.GetUserRequest{
		Id:    1,
		Token: "valid_token",
	}

	resp, err := server.GetUser(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.User)
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, "Cool", resp.User.FirstName)
	assert.Equal(t, "Kid", resp.User.LastName)
	assert.Equal(t, int32(12), resp.User.Age)
	assert.Equal(t, "valid_token", resp.User.Token)
}

func TestUpdateUser_InvalidUserData_UserNil_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.UpdateUserRequest{
		User: nil,
	}

	resp, err := server.UpdateUser(context.Background(), req)

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")

	statusErr, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")

	expectedErrorMsg := "Invalid user data"
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestUpdateUser_InvalidUserData_FirstNameEmpty_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "",
			LastName:  "Name",
			Age:       10,
		},
	}

	resp, err := server.UpdateUser(context.Background(), req)

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")

	statusErr, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")

	expectedErrorMsg := "Invalid user data"
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestUpdateUser_InvalidUserData_LastNameEmpty_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "User",
			LastName:  "",
			Age:       10,
		},
	}

	resp, err := server.UpdateUser(context.Background(), req)

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")

	statusErr, ok := status.FromError(err)
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")

	expectedErrorMsg := "Invalid user data"
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestUpdateUser_InvalidUserData_NegativeAge_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)

	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})

	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:        1,
			FirstName: "User",
			LastName:  "Name",
			Age:       -10,
		},
	}

	resp, err := server.UpdateUser(context.Background(), req)

	statusErr, ok := status.FromError(err)

	expectedErrorMsg := "Invalid user data"

	assert.Error(t, err, "Expected error for invalid user data")
	assert.Nil(t, resp, "Expected nil response for invalid user data")
	assert.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code(), "Expected InvalidArgument error")
	assert.Equal(t, expectedErrorMsg, statusErr.Message(), "Expected error message")

	mock.ExpectRollback()
}

func TestUpdateUser_InvalidToken_NotAuthenticated_ReturnsError(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "token"}).AddRow(1, "Cool", "Kid", 12, "valid_token")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	req := &pb.UpdateUserRequest{
		Id: 1,
		User: &pb.User{
			Id:        1,
			FirstName: "UpdatedFirstName",
			LastName:  "UpdatedLastName",
			Age:       10,
		},
		Token: "invalidToken",
	}

	resp, err := server.UpdateUser(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestUpdateUser_success(t *testing.T) {
	mockDB, mock, err := sqlmock.New()

	assert.Nil(t, err, "Failed to create mock DB: %v", err)
	defer mockDB.Close()

	dialector := postgres.New(postgres.Config{
		Conn:       mockDB,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.Nil(t, err, "Failed to open GORM DB: %v", err)

	server := &userServiceServer{DB: gormDB}

	mock.ExpectBegin()

	req := &pb.UpdateUserRequest{
		Id: 1,
		User: &pb.User{
			Id:        1,
			FirstName: "UpdatedFirstName",
			LastName:  "UpdatedLastName",
			Age:       10,
		},
		Token: "validToken",
	}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "token"}).AddRow("1", "FirstName", "LastName", 20, "validToken")
	newRow := sqlmock.NewRows([]string{"id", "first_name", "last_name", "age", "token"}).AddRow("1", "UpdatedFirstName", "UpdatedLastName", 10, "validToken")

	mock.ExpectQuery("SELECT").WithArgs(req.Id).WillReturnRows(rows)
	mock.ExpectQuery("UPDATE").WillReturnRows(newRow)
	mock.ExpectCommit()

	resp, err := server.UpdateUser(context.Background(), req)

	assert.NoError(t, err, "Unexpected error in CreatUser")
	assert.NotNil(t, resp, "Expected non-nil response")
	assert.Equal(t, "Created user successfully", resp.Message, "Unexpected response message")
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, "Cool", resp.User.FirstName)
	assert.Equal(t, "Kid", resp.User.LastName)
	assert.Equal(t, int32(10), resp.User.Age)
	assert.NotNil(t, resp.User.Token)
}
