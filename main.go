package main

import (
	"fmt"
	"os"
	"platform-tools/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
