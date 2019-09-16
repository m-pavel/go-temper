package tempern

import (
	"fmt"
	"testing"
)

func TestNative1(t *testing.T) {
	tm, err := New(0, 0, true)
	if err != nil {
		t.Fatal(err)
	}
	defer tm.Close()

	r, err := tm.Read()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(r)
}
