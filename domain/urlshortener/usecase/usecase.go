package usecase

import (
	"github.com/w-k-s/short-url/domain"
)

type Request interface{}

type Response interface{}

type UseCase interface {
	Execute(request Request) (Response, domain.Err)
}
