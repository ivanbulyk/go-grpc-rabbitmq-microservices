package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/ivanbulyk/go-grpc-rabbitmq-microservices/usermgmt"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	ID   int    `json:"id"`
}

const (
	address = "localhost: 50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var new_users = make(map[string]int32)

	new_users["Alice"] = 43
	new_users["Bob"] = 30

	connRabbit, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer connRabbit.Close()

	fmt.Println("Successfully connected to  RabbitMq instance")

	ch, err := connRabbit.Channel()
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Could not declare a queue: %v", err)
	}
	fmt.Println(q)

	for name, age := range new_users {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}

		log.Printf(`User Details:
NAME: %s
AGE: %d
ID: %d`, r.GetName(), r.GetAge(), r.GetId())
		userData := User{Name: r.GetName(), Age: int(r.GetAge()), ID: int(r.GetId())}
		body, err := json.Marshal(userData)
		if err != nil {
			log.Printf("Error encoding JSON %v", err)
		}

		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        body,
			},
		)
		if err != nil {
			log.Fatalf("Could not published message to a queue: %v", err)
		}
		fmt.Println("Succesfully Published Message to Queue")
	}
	params := &pb.GetUsersParams{}
	r, err := c.GetUsers(ctx, params)
	if err != nil {
		log.Fatalf("could not retrieve users: %v", err)
	}
	log.Print("\nUSER LIST: \n")
	fmt.Printf("r.GetUsers(): %v\n", r.GetUsers())

}
