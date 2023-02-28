package main

import (
	"context"

	"github.com/gin-gonic/gin"
	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetMovie(ctx context.Context, log zap.SugaredLogger, c *gin.Context, client moviesv1.MoviesServiceClient) {
	logger := log.With("method", "GetMovie")
	id := c.Query("id")
	logger.Infow("title", "title", id)

	if id == "" {
		logger.Errorw("title is required")
		c.JSON(400, gin.H{"error": "title is required"})
		return
	}

	moviesRequest := moviesv1.GetMovieByIdRequest{
		Id: id,
	}

	movie, err := client.GetMovieById(ctx, &moviesRequest)
	c.Writer.Header().Add("Content-Type", "application/json")
	if err != nil {
		errCode := status.Code(err)

		switch errCode {
		case codes.InvalidArgument:
			logger.Errorw("invalid argument", "error", err)
			c.JSON(400, gin.H{"error": "invalid argument"})
		case codes.NotFound:
			logger.Errorw("movie not found", "error", err)
			c.JSON(404, gin.H{"error": "movie not found"})
		default:
			logger.Errorw("internal server error", "error", err)
			c.JSON(500, gin.H{"error": "internal server error"})
		}
	} else {
		logger.Infow("movie", "movie", movie.Movie)
		c.JSON(200, gin.H{"movie": movie.Movie})
	}

}

type CreateMovieBody struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

func CreateMovie(ctx context.Context, log zap.SugaredLogger, c *gin.Context, client moviesv1.MoviesServiceClient) {
	logger := log.With("method", "CreateMovie")

	var body CreateMovieBody

	if err := c.BindJSON(&body); err != nil {
		logger.Errorw("invalid request body", "error", err)
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	if body.Title == "" {
		logger.Errorw("title is required")
		c.JSON(400, gin.H{"error": "title is required"})
		return
	}

	if body.Director == "" {
		logger.Errorw("director is required")
		c.JSON(400, gin.H{"error": "director is required"})
		return
	}

	if body.Year == 0 {
		logger.Errorw("year is required")
		c.JSON(400, gin.H{"error": "year is required"})
		return
	}

	moviesRequest := moviesv1.CreateMovieRequest{
		Movie: &moviesv1.Movie{
			Title:    body.Title,
			Director: body.Director,
			Year:     int32(body.Year),
		},
	}

	movie, err := client.CreateMovie(ctx, &moviesRequest)

	c.Writer.Header().Add("Content-Type", "application/json")

	if err != nil {
		errCode := status.Code(err)

		switch errCode {
		case codes.InvalidArgument:
			logger.Errorw("invalid argument", "error", err)
			c.JSON(400, gin.H{"error": "invalid argument"})
		default:
			logger.Errorw("internal server error", "error", err)
			c.JSON(500, gin.H{"error": "internal server error"})
		}
	} else {
		logger.Infow("movie", "movie", movie.Movie)
		c.JSON(200, gin.H{"movie": movie.Movie, "id": movie.Id})
	}

}

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
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		logger.Errorf("Failed to connect to server", "error", err)
		return err
	}

	defer conn.Close()

	moviesClient := moviesv1.NewMoviesServiceClient(conn)

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		GetMovie(ctx, log, ctx, moviesClient)
	})

	router.POST("/", func(ctx *gin.Context) {
		CreateMovie(ctx, log, ctx, moviesClient)
	})
	defer func() {
		err := router.Run(":8080")
		if err != nil {
			logger.Fatalw("Failed to run server", "error", err)
		}
	}()

	return nil
}
