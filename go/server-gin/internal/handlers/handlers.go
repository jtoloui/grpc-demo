package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"google.golang.org/grpc/codes"
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

	logger.Info("GetMovies", "page", page, "per_page", perPage)

	moviesRequest := moviesv1.GetMoviesRequest{
		Page:    int32(page),
		PerPage: int32(perPage),
	}

	movies, err := c.Client.GetMovies(c.Ctx, &moviesRequest)

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

	moviesRequest := moviesv1.CreateMovieRequest{
		Movie: &moviesv1.Movie{
			Title:    body.Title,
			Director: body.Director,
			Year:     int32(body.Year),
		},
	}

	movie, err := c.Client.CreateMovie(c.Ctx, &moviesRequest)

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
	logger.Infow("movie", "movie", movie.Movie)
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
	logger.Infow("id", "id", id)

	moviesRequest := moviesv1.GetMovieByIdRequest{
		Id: id,
	}

	movie, err := c.Client.GetMovieById(c.Ctx, &moviesRequest)
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
	logger.Infow("movie", "movie", movie.Movie)
	gc.JSON(200, gin.H{"movie": movie.Movie})
}
