package api

import (
	"encoding/json"
	"log"
	"lpp/lpp"
	"net/http"
	"strings"
)

func Run() {
	http.HandleFunc("/parse", handleParsing)
	log.Fatal(http.ListenAndServe(":1323", nil))
}

func handleParsing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var response string
	var data map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response = "unable to parse request body"
		http.Error(w, response, http.StatusBadRequest)
		return
	}

	source, exist := data["source"]
	if !exist {
		response = "source field is required"
		http.Error(w, response, http.StatusBadRequest)
		return
	}

	evaluated, erros := Parse(source.(string))
	w.WriteHeader(http.StatusOK)
	res := map[string]string{
		"evaluated": evaluated.Inspect(),
		"erros":     strings.Join(erros, " "),
	}

	json.NewEncoder(w).Encode(res)
}

func Parse(source string) (lpp.Object, []string) {
	lexer := lpp.NewLexer(source)
	parser := lpp.NewParser(lexer)
	env := lpp.NewEnviroment(nil)
	program := parser.ParseProgam()

	if len(parser.Errors()) > 0 {
		return nil, parser.Errors()
	}

	evaluated := lpp.Evaluate(program, env)
	return evaluated, nil
}
