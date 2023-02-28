package main

import (
	"context"
	"net"

	"github.com/jtoloui/grpc-demo/go/server-grpc/internal/config"
	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	moviesv1.UnimplementedMoviesServiceServer
	db  *mongo.Collection
	log *zap.SugaredLogger
}

type Movie struct {
	Title    string `json:"title"`
	Director string `json:"director"`
	Year     int    `json:"year"`
}

func (s *Server) GetMovieById(ctx context.Context, req *moviesv1.GetMovieByIdRequest) (*moviesv1.GetMovieByIdResponse, error) {
	logger := s.log.With("method", "GetMovieById")

	logger.Infow("GetMovieById", "id", req.Id)
	id, err := primitive.ObjectIDFromHex(req.Id)

	if err != nil {
		logger.Errorw("Error decoding bytes", "error", err)
		return nil, status.Error(codes.InvalidArgument, "Invalid ID")
	}
	filter := bson.D{primitive.E{Key: "_id", Value: primitive.ObjectID(id)}}

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
		},
	}, nil
}

func (s *Server) CreateMovie(ctx context.Context, req *moviesv1.CreateMovieRequest) (*moviesv1.CreateMovieResponse, error) {
	logger := s.log.With("method", "CreateMovie")

	logger.Infow("Creating movie", "movie", req.Movie)

	movie := Movie{
		Title:    req.Movie.Title,
		Director: req.Movie.Director,
		Year:     int(req.Movie.Year),
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

	moviesv1.RegisterMoviesServiceServer(s, &Server{
		db:  db,
		log: log.With("service", "MoviesService"),
	})

	if err := s.Serve(listener); err != nil {
		logger.Errorf("Failed to start server", "error", err)
		return err
	}
	logger.Infow("Server started", "port", config.GrpcPort)
	return nil
}
