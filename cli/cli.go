package cli

import (
	"fmt"

	"github.com/m-pavel/go-temper/pkg"
)

type Cli struct {
}

func (Cli) Main(t temper.Temper) {
	defer t.Close()

	r, err := t.Read()
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}
