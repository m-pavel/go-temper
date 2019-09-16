package cli

import (
	"fmt"

	"github.com/m-pavel/go-temper/pkg"
)

func main() {
	tm, err := temper.NewNative(0, 0, true)
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
