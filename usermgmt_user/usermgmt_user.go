package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/ivanbulyk/go-grpc-rabbitmq-microservices/usermgmt"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	conn *pgx.Conn
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf(" failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (s *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {

	log.Printf("Recieved: %v", in.GetName())
	createSql := `
	create table if not exists users(
		id SERIAL PRIMARY KEY,
		name TEXT,
		age INT
	);
	`

	_, err := s.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed: %v\n", err)
		os.Exit(1)
	}

	created_user := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	tx, err := s.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin failed: %v", err)
	}
	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1, $2)", created_user.Name, created_user.Age)
	if err != nil {
		log.Fatalf("tx.Exec failed: %v", err)
	}

	tx.Commit(context.Background())
	return created_user, nil
}

func (s *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {

	var users_list *pb.UserList = &pb.UserList{}
	rows, err := s.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users_list.Users = append(users_list.Users, &user)
	}

	return users_list, nil
}

func main() {
	database_url := "postgres://postgres:mysecretpassword@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), database_url)
	if err != nil {
		log.Fatalf("Unable to esteblish connection: %v", err)
	}
	defer conn.Close(context.Background())
	var user_mgmt_server *UserManagementServer = NewUserManagementServer()
	user_mgmt_server.conn = conn
	if err := user_mgmt_server.Run(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
