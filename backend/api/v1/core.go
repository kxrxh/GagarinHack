package v1

import (
	"github.com/gofiber/fiber/v2"

	"github.com/gagarin/backend/api/auth"
	"github.com/gagarin/backend/utils"
)

func SetupRoutesV1(v1 *fiber.Router) {
	v := (*v1).Group("/v1")

	auth.SetupAuth(&v)

	(*v1).Use(utils.Redirect)
}
