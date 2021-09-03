# go-grpc-rabbitmq-microservices

## An example of simple user management service  

   The service consists of two microservices: User and UserCash.  
   UserCashe can send a request to create a new user and User will create that user, store it in Postgres DB, take the created user datails from DB, and essentially send back to the UserCash, which in turn will save them to RabbitMQ Instance.  

Prerequisites

Postgres, either installed on your machine, or there is a Postgres instance that you can connect to.

The same for RabbitMQ

Go, any one of the three latest major releases of Go.

Protocol buffer compiler, protoc, version 3.  

Go plugins for the protocol compiler:  

Install the protocol compiler plugins for Go using the following commands:  

<p>$ go install google.golang.org/protobuf/cmd/protoc-gen-go  </p>
<p>$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc  </p>

<p>$ export PATH="$PATH:$(go env GOPATH)/bin"  </p>
Get the example code

Clone the repo:

$ git clone https://github.com/ivanbulyk/go-grpc-rabbitmq-microservices.git  

Change to the go-grpc-rabbitmq-microservices directory:

$ cd go-grpc-rabbitmq-microservices  

Run the example  
From the go-grpc-rabbitmq-microservices directory:

Compile and execute the User code:

$ go run usermgmt_user/usermgmt_user.go  

From a different terminal, compile and execute the UserCash code to see the UserCash output:  

$ go run usermgmt_user_cash/usermgmt_user_cash.go

If you can see something that is not an error, than..  Congratulations! You’ve just run a microservices interconnected application with gRPC.  
