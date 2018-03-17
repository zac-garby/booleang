package main

// This package defines the booleang (`bl`) command. Running `go install` in the
// /bl directory will generate a binary called `bl`.

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/Zac-Garby/booleang/lexer"
	"github.com/Zac-Garby/booleang/token"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func(c chan os.Signal) {
		<-c
		quit()
	}(c)

	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("booleang> ")
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}

		handleInput(line)
	}
}

func quit() {
	fmt.Println("quit")
	os.Exit(0)
}

func handleInput(input string) {
	if strings.TrimSpace(input) == "quit" {
		quit()
	}

	l := lexer.New(input, "repl")

	for tok := l(); tok.Type != token.EOF; tok = l() {
		fmt.Println(tok.String())
	}
}
