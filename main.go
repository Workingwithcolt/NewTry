package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"Body"`
}

func main() {

	todos := []Todo{}

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	app := fiber.New()
	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})
	//create Todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{} // stored the address of the Todo{}
		if err := c.BodyParser(todo); err != nil {
			return err
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error": "todo body is required"})
		}

		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		return c.Status(200).JSON(todo)

	})

	//update TODO
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		present := false
		for i, todo := range todos {
			if fmt.Sprint(todo.ID) == id {
				present = true
				todos = append(todos[:i], todos[i+1:]...)
				// the above code will do following
				//[:i] if index = 2 this will return 0 to 1 but not 2
				// [i+1:] that will return the array after 3 with 3 include
			}
		}
		if present {
			return c.Status(200).JSON(todos)
		} else {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found"})
		}
	})
	app.Listen(":" + PORT)
}
