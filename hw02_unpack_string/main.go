package main

import (
	"fmt"
	"log"
)

func main() {
	inputStrings := []string{
		/*
			"a4bc2d5e",
			"abcd",
			"3abc",
			"45",
			"aaa10b",
			"aaa0b",
			"",
			"d\n5abc",
		*/

		`qwe\4\5`,
		`qwe\45`,
		`qwe\\5`,
		`qw\ne`,
	}

	for _, inputString := range inputStrings {
		response, err := Unpack(inputString)
		if err != nil {
			log.Println(err)
		}

		fmt.Println(response)
	}
}
