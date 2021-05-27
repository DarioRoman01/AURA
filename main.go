package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"lpp/lpp"
	"os"
)

type CommandLine struct{}

func (cli *CommandLine) Repl() {
	lpp.StartRpl()
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage: ")
	fmt.Println("	file -path <path to your file> execute the given file")
	fmt.Println("	rpl - start the programing lenguage repl")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		os.Exit(0)
	}
}

func (cli *CommandLine) ReadFile(path string) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error: unable to read the give file")
	}

	if len(content) == 0 {
		return
	}

	lexer := lpp.NewLexer(string(content))
	parser := lpp.NewParser(lexer)
	env := lpp.NewEnviroment(nil)
	program := parser.ParseProgam()

	if len(parser.Errors()) > 0 {
		for _, err := range parser.Errors() {
			fmt.Println(err)
		}
	}

	evaluated := lpp.Evaluate(program, env)
	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}

func (cli *CommandLine) Start() {
	cli.ValidateArgs()

	fileCmd := flag.NewFlagSet("file", flag.ExitOnError)
	rplCmd := flag.NewFlagSet("rpl", flag.ExitOnError)

	filePath := fileCmd.String("path", "", "the path you want to execute")

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
