package services

import (
	"context"
	"echo-react-serve/server/repositories"
)

type Service interface {
	DoSomething(ctx context.Context)
}

type service struct {
	r repositories.Repo
}

func NewService(r repositories.Repo) Service {
	return &service{r}
}

func (s *service) DoSomething(ctx context.Context) {}
