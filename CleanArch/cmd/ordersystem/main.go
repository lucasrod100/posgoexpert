package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/lucasrod100/posgoexpert/CleanArch/configs"
	"github.com/lucasrod100/posgoexpert/CleanArch/internal/event/handler"
	"github.com/lucasrod100/posgoexpert/CleanArch/internal/infra/graph"
	"github.com/lucasrod100/posgoexpert/CleanArch/internal/infra/grpc/pb"
	"github.com/lucasrod100/posgoexpert/CleanArch/internal/infra/grpc/service"
	"github.com/lucasrod100/posgoexpert/CleanArch/internal/infra/web/webserver"
	"github.com/lucasrod100/posgoexpert/CleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(configs.RabbitMQUser, configs.RabbitMQPassword, configs.RabbitMQPort)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	newWebServer := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	newWebServer.AddHandler(webserver.MethodHandler{
		Path:    "/order",
		Method:  "POST",
		Handler: webOrderHandler.Create,
	})
	newWebServer.AddHandler(webserver.MethodHandler{
		Path:    "/order",
		Method:  "GET",
		Handler: webOrderHandler.List,
	})
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go newWebServer.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel(user, pass, port string) *amqp.Channel {
	conn, err := amqp.Dial("amqp://" + user + ":" + pass + "@rabbitmq:" + port + "/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
