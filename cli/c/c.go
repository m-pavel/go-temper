package main

import (
	"github.com/m-pavel/go-temper/cli"
)

func main() {
	tm, err := temperc.New(0, 0, true)
	if err != nil {
		panic(err)
	}
	cli.Cli{}.Main(tm)
}
