package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"katan/src"
	"os"
)

type CommandLine struct{}

func (cli *CommandLine) Repl() {
	src.StartRpl()
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage: ")
	fmt.Println("	file -path <path to your file> - will execute the given file")
	fmt.Println("	rpl - Starts the repl")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(0)
	}
}

func (cli *CommandLine) ReadFile(path string) {
	source, err := ioutil.ReadFile(path)
	if err != nil {
		cli.PrintUsage()
		os.Exit(0)
	}

	if len(source) == 0 {
		return
	}

	lexer := src.NewLexer(string(source))
	parser := src.NewParser(lexer)
	env := src.NewEnviroment(nil)
	program := parser.ParseProgam()

	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Println(err)
		}

		os.Exit(0)
	}

	evaluated := src.Evaluate(program, env)
	if evaluated != nil && evaluated != src.SingletonNUll {
		fmt.Println(evaluated.Inspect())
	}
}

func (cli *CommandLine) Start() {
	cli.ValidateArgs()

	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	rplCmd := flag.NewFlagSet("rpl", flag.ExitOnError)

	filePath := fileCmd.String("path", "", "the path to your file")

	switch os.Args[1] {
	case "file":
		err := fileCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

	case "rpl":
		err := rplCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	if fileCmd.Parsed() {
		if *filePath == "" {
			cli.PrintUsage()
			os.Exit(0)
		}

		cli.ReadFile(*filePath)
	}

	if rplCmd.Parsed() {
		cli.Repl()
	}
}

func main() {
	cli := CommandLine{}
	cli.Start()
}
