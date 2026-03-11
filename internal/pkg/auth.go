package pkg

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetUserId(c *fiber.Ctx) (uint64, bool) {
	userIdStr, ok := c.Locals("userId").(string)
	if !ok {
		return 0, false
	}
	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	return userId, err == nil
}
