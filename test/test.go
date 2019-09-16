package main

import "fmt"

func main() {
	var buf0 uint16
	var buf1 uint16
	buf0 = 25
	buf1 = 165
	temperature := (buf1 & 0xFF) + (buf0 << 8)
	fmt.Println(buf1 & 0xFF)
	fmt.Println(buf0 << 8)
	fmt.Println(temperature)
	temp := -39.7 + .01*float64(temperature)
	fmt.Printf("%d %f\n", temperature, temp)
}
