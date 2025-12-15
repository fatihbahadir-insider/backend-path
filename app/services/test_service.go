package services

import (
	"backend-path/app/repository"
	"backend-path/utils"

	"github.com/gofiber/fiber/v2"
)

type TestService struct {}

func (s *TestService) newsRepo() *repository.TestRepository {
	return  new(repository.TestRepository)
}

func (s *TestService) List(ctx *fiber.Ctx) error {
	tests, err := s.newsRepo().GetAll()
	if err != nil {
		return utils.JsonErrorInternal(ctx, err, "E_TEST_LIST")
	}

	return utils.JsonSuccess(ctx, tests)
}