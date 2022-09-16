package rest

import "github.com/gofiber/fiber/v2"

func GlobalHandler(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).SendString(err.Error())
}
