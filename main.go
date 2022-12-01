package main

import (
	e "aura/src/evaluator"
	l "aura/src/lexer"
	obj "aura/src/object"
	p "aura/src/parser"
	"encoding/json"
	"html"
	"log"
	"net/http"
)

func handleEvaluation(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	source, ok := body["source"]
	if !ok {
		http.Error(w, "no source provided", http.StatusBadRequest)
		return
	}

	lexer := l.NewLexer(source)
	parser := p.NewParser(lexer)
	env := obj.NewEnviroment(nil)
	program := parser.ParseProgam()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if len(parser.Errors()) > 0 {
		json.NewEncoder(w).Encode(map[string][]string{
			"errors": parser.Errors(),
		})
		// we dont evaluate the program if it has syntax errors
		return
	}

	evaluated := e.Evaluate(program, env)
	if evaluated != nil && evaluated != obj.SingletonNUll {
		json.NewEncoder(w).Encode(map[string]string{
			"result": html.EscapeString(evaluated.Inspect()),
		})
	}
}

func main() {
	http.HandleFunc("/", handleEvaluation)
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
