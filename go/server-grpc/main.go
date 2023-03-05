package main

import (
	"context"
	"net"

	"github.com/jtoloui/grpc-demo/go/server-grpc/internal/config"
	"github.com/jtoloui/grpc-demo/go/server-grpc/internal/service"
	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Create a new grpc server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log, _ := zap.NewProduction()
	defer log.Sync()

	logger := log.Sugar()

	if err := run(ctx, *logger); err != nil {
		logger.Fatalw("Failed to run server", "error", err)
	}
}

func run(ctx context.Context, log zap.SugaredLogger) error {
	// setup mongo db connection
	logger := log.With("method", "run")
	defer logger.Sync()

	config := config.GetConfig()
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	logger.Info("Connecting to MongoDB")
	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			logger.Fatalw("Failed to disconnect MongoDB", "error", err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		logger.Errorf("Failed to ping MongoDB", "error", err)
		return err
	}

	logger.Info("Connected to MongoDB")

	// setup grpc server

	listener, err := net.Listen("tcp", config.GrpcPort)
	if err != nil {
		logger.Errorf("Failed to listen", "error", err)
		return err
	}

	s := grpc.NewServer()

	db := client.Database("grpc-demo").Collection("movies")

	service := service.NewService(db, log.With("service", "MoviesService"))

	moviesv1.RegisterMoviesServiceServer(s, service)

	if err := s.Serve(listener); err != nil {
		logger.Errorf("Failed to start server", "error", err)
		return err
	}
	logger.Infow("Server started", "port", config.GrpcPort)
	return nil
}
