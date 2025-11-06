package env

import "gorm.io/gorm"

type Env struct {
	database *gorm.DB
}

type Option func(*Env) *Env

func NewEnv(opts ...Option) *Env {
	env := &Env{}
	for _, opt := range opts {
		env = opt(env)
	}
	return env
}

func (e *Env) Database() *gorm.DB {
	return e.database
}

func (e *Env) Close() error {
	if e.database != nil {
		ins, _ := e.database.DB()
		ins.Close() //nolint
	}
	return nil
}

func WithDatabase(db *gorm.DB) Option {
	return func(env *Env) *Env {
		env.database = db
		return env
	}
}
