package presentation

import (
	"clean_architect/application/usecase"
	"clean_architect/config"
	"clean_architect/env"
	"clean_architect/presentation/http"
	"context"
)

type Presenter interface {
	Start(ctx context.Context) error
	AppHttp() http.AppHttp
}

type presenter struct {
	appHttp http.AppHttp
}

func NewPresenter(cfg *config.Config, env *env.Env) Presenter {
	usecase := usecase.NewUsecase()
	appHttp := http.New(usecase)
	return &presenter{appHttp: appHttp}
}

func (p *presenter) Start(ctx context.Context) error {
	p.appHttp.Routes(ctx)
	return nil
}

func (p *presenter) AppHttp() http.AppHttp {
	return p.appHttp
}
