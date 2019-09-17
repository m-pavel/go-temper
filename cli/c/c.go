package main

import (
	"fmt"

	"github.com/m-pavel/go-temper/pkg-c"
)

func main() {
	tm, err := temperc.New(0, 0, true)
	if err != nil {
		panic(err)
	}
	defer tm.Close()

	r, err := tm.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)

}
