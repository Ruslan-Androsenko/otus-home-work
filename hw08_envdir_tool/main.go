package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		environments, err := ReadDir(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}

		for key, value := range environments {
			fmt.Printf("%s is (%s)\n", key, value.Value)
		}
	}
}
