package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) > 2 {
		environments, err := ReadDir(os.Args[1])
		if err != nil {
			log.Fatalln(err)
		}

		RunCmd(os.Args[2:], environments)
	}
}

// args := os.Args
// fmt.Sprint(args)
// fmt.Println(environments)
// ./testdata/env /bin/bash arg1=1 arg2=2

//  result='[
// ./go-envdir
// /home/ruslan/workspace/projects/golang/otus-home-work-main/hw08_envdir_tool/testdata/env
// /bin/bash
// /home/ruslan/workspace/projects/golang/otus-home-work-main/hw08_envdir_tool/testdata/echo.sh
// arg1=1
// arg2=2
// ]'
