package handlers

import (
	"context"

	moviesv1 "github.com/jtoloui/proto-store/go/movies/v1"
	"go.uber.org/zap"
)

type Config struct {
	Log    *zap.SugaredLogger
	Ctx    context.Context
	Client moviesv1.MoviesServiceClient
}

func NewConfig(log *zap.SugaredLogger, ctx context.Context, client moviesv1.MoviesServiceClient) *Config {
	return &Config{
		Log:    log,
		Ctx:    ctx,
		Client: client,
	}
}
