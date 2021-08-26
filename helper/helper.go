package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/online-store/entity"
)

func WriteResponse(c *fiber.Ctx, response entity.Response) error {
	return c.Status(int(response.Code)).JSON(response)
}
