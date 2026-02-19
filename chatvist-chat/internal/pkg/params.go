package pkg

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateParamsId(c *fiber.Ctx) (uint64, error) {
	idStr := c.Params("id")
	if len(strings.TrimSpace(idStr)) == 0 {
		return 0, errors.New("ID de producto no proporcionado")
	}

	id, newErr := strconv.ParseUint(idStr, 10, 64)
	if newErr != nil {
		return 0, errors.New("ID de producto inválido: " + newErr.Error())
	}

	return id, nil
}
