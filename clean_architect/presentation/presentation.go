package presentation

import (
	"clean_architect/application/usecase"
	"clean_architect/config"
	"clean_architect/env"
	"clean_architect/infrastructure/jwt"
	"clean_architect/infrastructure/persistent/repository"
	"clean_architect/presentation/http"
	"context"
	"time"
)

type Presenter interface {
	Start(ctx context.Context) error
	AppHttp() http.AppHttp
}

type presenter struct {
	appHttp http.AppHttp
}

func NewPresenter(cfg *config.Config, env *env.Env) Presenter {
	// Initialize JWT service
	jwtService := jwt.NewJWTService(
		cfg.JWT.SecretKey,
		time.Duration(cfg.JWT.TokenLifetime)*time.Hour,
	)

	// Initialize repositories
	repos := repository.NewRepos(env, env.Database())

	// Initialize use cases
	usecases := usecase.NewUsecase(repos, jwtService)

	// Initialize HTTP layer
	appHttp := http.New(usecases)
	return &presenter{appHttp: appHttp}
}

func (p *presenter) Start(ctx context.Context) error {
	p.appHttp.Routes(ctx)
	return nil
}

func (p *presenter) AppHttp() http.AppHttp {
	return p.appHttp
}
