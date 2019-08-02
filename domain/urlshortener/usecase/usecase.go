package usecase

import(
	"github.com/w-k-s/short-url/domain"
)

interface Request{}

interface Response{}

interface UseCase{
	Execute(request Request) (response Response, domain.Err)
}