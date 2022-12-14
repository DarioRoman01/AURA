package main

import (
	e "aura/src/evaluator"
	l "aura/src/lexer"
	obj "aura/src/object"
	p "aura/src/parser"
	"fmt"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Request struct {
	Source string `json:"source"`
}

func handleEvaluation2(c *fiber.Ctx) error {
	var body Request
	if err := c.BodyParser(&body); err != nil {
		fmt.Println(err)
		return err
	}

	source := body.Source
	lexer := l.NewLexer(source)
	parser := p.NewParser(lexer)
	env := obj.NewEnviroment(nil)
	program := parser.ParseProgam()

	if len(parser.Errors()) > 0 {
		return c.Status(http.StatusOK).JSON(map[string][]string{
			"errors": parser.Errors(),
		})
	}

	evaluated := e.Evaluate(program, env)
	if evaluated != nil && evaluated != obj.SingletonNUll {
		return c.Status(http.StatusOK).JSON(map[string]string{
			"result": evaluated.Inspect(),
		})
	}

	return c.Status(http.StatusOK).JSON(map[string]string{
		"result": "no content",
	})
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Post("/", handleEvaluation2)
	log.Fatalf("Error: %s", app.Listen(":8080"))
}
