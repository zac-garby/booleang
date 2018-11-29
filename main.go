package main

// This package defines the booleang (`bl`) command. Running `go install` in the
// /bl directory will generate a binary called `bl`.

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/zac-garby/booleang/parser"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(c chan os.Signal) {
		<-c
		quit()
	}(c)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("no file specified...\nexecute a file by passing it's path as an argument")
		os.Exit(1)
	}

	file, err := os.Open(args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	handleFile(args[1], file)
}

func quit() {
	fmt.Println("quit")
	os.Exit(0)
}

func handleFile(path string, file *os.File) {
	text, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	filename := filepath.Base(path)

	p := parser.New(string(text), filename)
	prog, err := p.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(prog)
}
