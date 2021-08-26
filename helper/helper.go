package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/online-store/entity"
)

// fungsi untuk mengembalikan response
func WriteResponse(c *fiber.Ctx, response entity.Response) error {
	return c.Status(int(response.Code)).JSON(response)
}
