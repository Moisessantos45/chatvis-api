package pkg

import "github.com/gofiber/fiber/v2"

func ResponseJson(c *fiber.Ctx, status int, message string, errorCode string, details any) error {
	response := fiber.Map{
		"success": status >= 200 && status < 300,
		"message": message,
		"status":  status,
	}

	if errorCode != "" {
		response["error"] = fiber.Map{
			"code":    errorCode,
			"details": details,
		}
	} else {
		response["data"] = details
	}

	return c.Status(status).JSON(response)
}
