package api

import (
	"encoding/json"

	"github.com/gofiber/fiber/v3"

	"github.com/kxrxh/kvdb/storage"
)

func setupV1Router(app *fiber.Router) {
	(*app).Get("/:id", getHandler)
	(*app).Post("/", postHandler)
	(*app).Delete("/:id", deleteHandler)
	(*app).Delete("prefix/:id", deletePrefixHandler)

}

// getHandler is a handler for GET requests
func getHandler(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing id",
		})
	}

	value, found := storage.KeyValueStore.Get(id)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"key":    id,
		"value":  value,
		"exists": found,
	})
}

type postRequestBody struct {
	Key string `json:"key"`
	Val string `json:"value"`
}

func postHandler(c fiber.Ctx) error {
	c.Accepts("application/json")
	body := new(postRequestBody)

	// !Warning: I was unable to use `c.BodyParser(body)` because I don't have one in my v3. So. I'm using `c.Body()`
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if body.Key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing key"})
	}

	storage.KeyValueStore.Set(body.Key, body.Val)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"key":   body.Key,
		"value": body.Val,
	})
}

func deleteHandler(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing id",
		})
	}

	err := storage.KeyValueStore.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"key": id,
	})
}

func deletePrefixHandler(c fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "missing id",
		})
	}

	err := storage.KeyValueStore.DeleteSubtree(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"key": id,
	})
}
