package main

import (
	"fmt"
	"os"
)

func main() {
	err := Run()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
