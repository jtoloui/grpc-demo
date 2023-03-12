package service

import (
	"context"

	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Service struct {
	moviesv1.UnimplementedMoviesServiceServer
	db  *mongo.Collection
	log *zap.SugaredLogger
}

type Movie struct {
	Title    string             `json:"title"`
	Director string             `json:"director"`
	Year     int                `json:"year"`
	ID       primitive.ObjectID `bson:"_id"`
}

func NewService(db *mongo.Collection, log *zap.SugaredLogger) *Service {
	return &Service{
		db:  db,
		log: log,
	}
}

func (s *Service) GetMovies(ctx context.Context, req *moviesv1.GetMoviesRequest) (*moviesv1.GetMoviesResponse, error) {

	logger := s.log.With("method", "GetMovies")
	md, ok := metadata.FromIncomingContext(ctx)

	var tracerId string
	if !ok {
		tracerId = "no-tracer-id"
	} else {
		tracerId = md.Get("x-tracer-id")[0]
	}
	logger.Infow("inbound request", "x-tracer-id", tracerId)

	page := int64(req.Page)
	perPage := int64(req.PerPage)

	var movies []Movie
	opts := options.Find().SetSkip((page - 1) * perPage).SetLimit(perPage)

	cursor, err := s.db.Find(ctx, bson.D{}, opts)

	if err != nil {
		logger.Errorw("Error decoding bytes", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if err = cursor.All(ctx, &movies); err != nil {
		logger.Errorw("Error decoding bytes", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	var moviesResponse []*moviesv1.Movie
	for _, movie := range movies {
		moviesResponse = append(moviesResponse, &moviesv1.Movie{
			Title:    movie.Title,
			Director: movie.Director,
			Year:     int32(movie.Year),
			Id:       movie.ID.Hex(),
		})
	}

	est, estErr := s.db.EstimatedDocumentCount(ctx)

	if estErr != nil {
		logger.Errorw("Error decoding bytes", "error", estErr)
		return nil, status.Error(codes.Internal, "finding estimate failed")
	}

	return &moviesv1.GetMoviesResponse{
		Movies: moviesResponse,
		Total:  int32(est),
	}, nil

}

func (s *Service) GetMovieById(ctx context.Context, req *moviesv1.GetMovieByIdRequest) (*moviesv1.GetMovieByIdResponse, error) {
	logger := s.log.With("method", "GetMovieById")

	if req.Id == "" {
		logger.Errorw("Error decoding bytes", "error", "id is empty")
		return nil, status.Error(codes.InvalidArgument, "Invalid ID")
	}

	md, ok := metadata.FromIncomingContext(ctx)

	var tracerId string
	if !ok {
		tracerId = "no-tracer-id"
	} else {
		tracerId = md.Get("x-tracer-id")[0]
	}

	var id primitive.ObjectID

	if err := id.UnmarshalText([]byte(req.Id)); err != nil {
		logger.Errorw("Error decoding bytes", "error", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid ID")
	}

	logger.Infow("inbound request", "x-tracer-id", tracerId, "id", id.Hex())

	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	var movie Movie
	findErr := s.db.FindOne(ctx, filter).Decode(&movie)

	if findErr != nil {
		logger.Errorw("Error decoding bytes", "error", findErr)
		return nil, status.Error(codes.NotFound, "Movie not found")
	}

	return &moviesv1.GetMovieByIdResponse{
		Movie: &moviesv1.Movie{
			Title:    movie.Title,
			Director: movie.Director,
			Year:     int32(movie.Year),
			Id:       movie.ID.Hex(),
		},
	}, nil
}

func (s *Service) CreateMovie(ctx context.Context, req *moviesv1.CreateMovieRequest) (*moviesv1.CreateMovieResponse, error) {
	logger := s.log.With("method", "CreateMovie")

	md, ok := metadata.FromIncomingContext(ctx)

	var tracerId string
	if !ok {
		tracerId = "no-tracer-id"
	} else {
		tracerId = md.Get("x-tracer-id")[0]
	}

	logger.Infow("Creating movie", "movie", req.Movie, "x-tracer-id", tracerId)

	movie := Movie{
		Title:    req.Movie.Title,
		Director: req.Movie.Director,
		Year:     int(req.Movie.Year),
		ID:       primitive.NewObjectID(),
	}

	res, err := s.db.InsertOne(ctx, movie)

	if err != nil {
		logger.Errorw("Error inserting movie", "error", err)
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &moviesv1.CreateMovieResponse{
		Id: res.InsertedID.(primitive.ObjectID).Hex(),
		Movie: &moviesv1.Movie{
			Title:    movie.Title,
			Director: movie.Director,
			Year:     int32(movie.Year),
		},
	}, nil
}
