package repository

import (
	"clean_architect/env"
	"context"
	"log"

	"gorm.io/gorm"
)

type Repos interface {
	UserRepository() UserRepository
	WithTransaction(ctx context.Context, fn func(Repos) error) (err error)
}

type repos struct {
	env *env.Env
	db  *gorm.DB

	userRepository UserRepository
}

func NewRepos(env *env.Env, db *gorm.DB) Repos {
	userRepository := NewUserRepository(env.Database())
	return &repos{
		env: env,
		db:  db,

		userRepository: userRepository,
	}
}

func (r *repos) UserRepository() UserRepository {
	return r.userRepository
}

func (r *repos) WithTransaction(ctx context.Context, fn func(Repos) error) (err error) {
	tx := r.db.Begin()
	tr := NewRepos(r.env, tx)

	err = tx.Error
	if err != nil {
		return
	}

	defer func() {
		if p := recover(); p != nil { // nolint
			// a panic occurred, rollback and repanic
			log.Printf("ERROR: WithTransaction - panic: %v", p)
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			log.Printf("ERROR: WithTransaction - error: %v", err)
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit().Error
		}
	}()

	err = fn(tr)

	return err
}
