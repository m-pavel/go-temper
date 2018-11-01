package main

import (
	"github.com/m-pavel/go-temper/pkg"
	"log"
)

func main() {
	t := temper.New(0, 5)
	if t == nil {
		log.Fatal("Unable to create device")
	}
	res,err  := t.Read()
	log.Println(res)
	log.Println(err)
	t.Close()
}
