package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GetMovies is a handler for getting movies list
func (c *Config) GetMovies(gc *gin.Context) {
	logger := c.Log.With("method", "GetMovies")

	pageQuery := gc.DefaultQuery("page", "1")
	perPageQuery := gc.DefaultQuery("per_page", "10")

	page, err := strconv.Atoi(pageQuery)
	if err != nil {
		logger.Errorw("invalid page", "error", err)
		gc.JSON(400, gin.H{"error": "invalid page"})
		return
	}

	perPage, err := strconv.Atoi(perPageQuery)
	if err != nil {
		logger.Errorw("invalid per_page", "error", err)
		gc.JSON(400, gin.H{"error": "invalid per_page"})
		return
	}

	tracerId, ok := gc.Get("X-Tracer-Id")
	if !ok {
		logger.Errorw("tracer id not found")
		tracerId = "unknown"
	}
	md := metadata.Pairs("X-Tracer-Id", tracerId.(string))

	ctx := metadata.NewOutgoingContext(c.Ctx, md)

	logger.Infow("GetMovies", "page", page, "per_page", perPage, "x-tracer-id", tracerId)

	moviesRequest := moviesv1.GetMoviesRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	}

	movies, err := c.Client.GetMovies(ctx, &moviesRequest, grpc.Header(&md))

	gc.Writer.Header().Add("Content-Type", "application/json")

	if err != nil {
		errCode := status.Code(err)

		switch errCode {
		case codes.InvalidArgument:
			logger.Errorw("invalid argument", "error", err)
			gc.JSON(400, gin.H{"error": "invalid argument"})
			return
		default:
			logger.Errorw("internal server error", "error", err)
			gc.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}

	movieList := make([]moviesv1.Movie, 0, perPage)

	for _, movie := range movies.Movies {
		movieList = append(movieList, *movie)
	}
	gc.JSON(200, gin.H{"movies": movieList, "total": movies.Total})

}

// CreateMovies is a handler for creating a movie
type CreateMovieBody struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

func (c *Config) CreateMovie(gc *gin.Context) {
	logger := c.Log.With("method", "CreateMovie")

	var body CreateMovieBody

	if err := gc.BindJSON(&body); err != nil {
		logger.Errorw("invalid request body", "error", err)
		gc.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	if body.Title == "" {
		logger.Errorw("title is required")
		gc.JSON(400, gin.H{"error": "title is required"})
		return
	}

	if body.Director == "" {
		logger.Errorw("director is required")
		gc.JSON(400, gin.H{"error": "director is required"})
		return
	}

	if body.Year == 0 {
		logger.Errorw("year is required")
		gc.JSON(400, gin.H{"error": "year is required"})
		return
	}

	tracerId, ok := gc.Get("X-Tracer-Id")
	if !ok {
		logger.Errorw("tracer id not found")
		tracerId = "unknown"
	}
	md := metadata.Pairs("X-Tracer-Id", tracerId.(string))

	ctx := metadata.NewOutgoingContext(c.Ctx, md)

	moviesRequest := moviesv1.CreateMovieRequest{
		Movie: &moviesv1.Movie{
			Title:    body.Title,
			Director: body.Director,
			Year:     int32(body.Year),
		},
	}

	movie, err := c.Client.CreateMovie(ctx, &moviesRequest)

	gc.Writer.Header().Add("Content-Type", "application/json")

	if err != nil {
		errCode := status.Code(err)

		switch errCode {
		case codes.InvalidArgument:
			logger.Errorw("invalid argument", "error", err)
			gc.JSON(400, gin.H{"error": "invalid argument"})
			return
		default:
			logger.Errorw("internal server error", "error", err)
			gc.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}
	logger.Infow("CreateMovie", "movie", movie.Movie, "id", movie.Id, "x-tracer-id", tracerId)
	gc.JSON(200, gin.H{"movie": movie.Movie, "id": movie.Id})

}

// GetMovieByID is a handler for getting a movie by id
func (c *Config) GetMovieById(gc *gin.Context) {
	logger := c.Log.With("method", "GetMovie")
	id := gc.Param("id")

	if id == "" {
		logger.Errorw("id is required")
		gc.JSON(400, gin.H{"error": "id is required"})
		return
	}
	tracerId, ok := gc.Get("X-Tracer-Id")
	if !ok {
		logger.Errorw("tracer id not found")
		tracerId = "unknown"
	}
	md := metadata.Pairs("X-Tracer-Id", tracerId.(string))

	ctx := metadata.NewOutgoingContext(c.Ctx, md)

	logger.Infow("GetMovieById", "id", id, "x-tracer-id", tracerId)

	moviesRequest := moviesv1.GetMovieByIdRequest{
		Id: id,
	}

	movie, err := c.Client.GetMovieById(ctx, &moviesRequest)
	gc.Writer.Header().Add("Content-Type", "application/json")
	if err != nil {
		errCode := status.Code(err)

		switch errCode {
		case codes.InvalidArgument:
			logger.Errorw("invalid argument", "error", err)
			gc.JSON(400, gin.H{"error": "invalid argument"})
			return
		case codes.NotFound:
			logger.Errorw("movie not found", "error", err)
			gc.JSON(404, gin.H{"error": "movie not found"})
			return
		default:
			logger.Errorw("internal server error", "error", err)
			gc.JSON(500, gin.H{"error": "internal server error"})
			return
		}
	}
	gc.JSON(200, gin.H{"movie": movie.Movie})
}
