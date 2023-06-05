package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	hello := "Hello, OTUS!"
	rhello := stringutil.Reverse(hello)
	fmt.Println(rhello)
}
