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

// func handleEvaluation(w http.ResponseWriter, r *http.Request) {
// 	var body map[string]string
// 	err := json.NewDecoder(r.Body).Decode(&body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	source, ok := body["source"]
// 	if !ok {
// 		http.Error(w, "no source provided", http.StatusBadRequest)
// 		return
// 	}

// 	lexer := l.NewLexer(source)
// 	parser := p.NewParser(lexer)
// 	env := obj.NewEnviroment(nil)
// 	program := parser.ParseProgam()

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

// 	if len(parser.Errors()) > 0 {
// 		json.NewEncoder(w).Encode(map[string][]string{
// 			"errors": parser.Errors(),
// 		})
// 		// we dont evaluate the program if it has syntax errors
// 		return
// 	}

// 	evaluated := e.Evaluate(program, env)
// 	if evaluated != nil && evaluated != obj.SingletonNUll {
// 		json.NewEncoder(w).Encode(map[string]string{
// 			"result": evaluated.Inspect(),
// 		})
// 	}
// }

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
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Post("/", handleEvaluation2)
	log.Fatalf("Error: %s", app.Listen(":8080"))
}
