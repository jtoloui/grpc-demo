package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jtoloui/grpc-demo/go/server-gin/internal/handlers"
	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	logger := log.With("method", "run")
	defer logger.Sync()

	// create gin server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		logger.Errorf("Failed to connect to server", "error", err)
		return err
	}

	defer conn.Close()

	moviesClient := moviesv1.NewMoviesServiceClient(conn)

	router := gin.Default()

	handlers := handlers.NewConfig(&log, ctx, moviesClient)

	router.GET("/", handlers.GetMovies)

	router.GET("/:id", handlers.GetMovieById)

	router.POST("/", handlers.CreateMovie)

	defer func() {
		err := router.Run(":8080")
		if err != nil {
			logger.Fatalw("Failed to run server", "error", err)
		}
	}()

	return nil
}
