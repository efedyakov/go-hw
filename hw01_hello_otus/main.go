package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	hello := stringutil.Reverse("Hello, OTUS!")
	fmt.Println(hello)
}
