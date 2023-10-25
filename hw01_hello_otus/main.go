package main

import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

func main() {
	sourceMessage := "Hello, OTUS!"

	// Использовал этот пакет т.к. 31.07.2023 он был переименован,
	// Commits on Jul 31, 2023: Renamed from stringutil/reverse.go to hello/reverse/reverse.go
	// see: https://github.com/golang/example/commits/master/hello/reverse/reverse.go
	reverseMessage := reverse.String(sourceMessage)

	fmt.Println(reverseMessage)
}
