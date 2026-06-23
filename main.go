package main

import (
	"fmt"
	"git_lite/commands"
)

func main() {
	err := commands.Log()
	if err != nil {
		fmt.Println(err)
	}
}