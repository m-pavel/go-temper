package main

import (
	"fmt"
)

func main() {
	tm, err := tempern.New(0, 0, true)
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
