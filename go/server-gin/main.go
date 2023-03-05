package main

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GetMovieById(ctx context.Context, log zap.SugaredLogger, c *gin.Context, client moviesv1.MoviesServiceClient) {
	logger := log.With("method", "GetMovie")
	id := c.Param("id")

	if id == "" {
		logger.Errorw("id is required")
		c.JSON(400, gin.H{"error": "id is required"})
		return
	}
	logger.Infow("id", "id", id)

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
			return
		case codes.NotFound:
			logger.Errorw("movie not found", "error", err)
			c.JSON(404, gin.H{"error": "movie not found"})
			return
		default:
			logger.Errorw("internal server error", "error", err)
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}
	logger.Infow("movie", "movie", movie.Movie)
	c.JSON(200, gin.H{"movie": movie.Movie})
	return

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
			return
		default:
			logger.Errorw("internal server error", "error", err)
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}
	logger.Infow("movie", "movie", movie.Movie)
	c.JSON(200, gin.H{"movie": movie.Movie, "id": movie.Id})

}

func GetMovies(ctx context.Context, log zap.SugaredLogger, c *gin.Context, client moviesv1.MoviesServiceClient) {
	logger := log.With("method", "GetMovies")

	pageQuery := c.DefaultQuery("page", "1")
	perPageQuery := c.DefaultQuery("per_page", "10")

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		logger.Errorw("invalid page", "error", err)
		c.JSON(400, gin.H{"error": "invalid page"})
		return
	}

	perPage, err := strconv.Atoi(perPageQuery)
	if err != nil {
		logger.Errorw("invalid per_page", "error", err)
		c.JSON(400, gin.H{"error": "invalid per_page"})
		return
	}

	logger.Info("GetMovies", "page", page, "per_page", perPage)

	moviesRequest := moviesv1.GetMoviesRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	}

	movies, err := client.GetMovies(ctx, &moviesRequest)

	c.Writer.Header().Add("Content-Type", "application/json")

	if err != nil {
		errCode := status.Code(err)

		switch errCode {
		case codes.InvalidArgument:
			logger.Errorw("invalid argument", "error", err)
			c.JSON(400, gin.H{"error": "invalid argument"})
			return
		default:
			logger.Errorw("internal server error", "error", err)
			c.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}

	movieList := make([]moviesv1.Movie, 0, perPage)

	for _, movie := range movies.Movies {
		movieList = append(movieList, *movie)
	}
	c.JSON(200, gin.H{"movies": movieList, "total": movies.Total})

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
		GetMovies(ctx, log, ctx, moviesClient)
	})

	router.GET("/:id", func(ctx *gin.Context) {
		GetMovieById(ctx, log, ctx, moviesClient)
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
