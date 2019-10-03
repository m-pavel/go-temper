package main

import (
	"github.com/m-pavel/go-temper/cli"
	"github.com/m-pavel/go-temper/pkg-native"
)

func main() {
	tm, err := tempern.New(0, 0, true)
	if err != nil {
		panic(err)
	}
	cli.Cli{}.Main(tm)
}
