package main

import (
	"fmt"
	"git_lite/commands"
)

func main() {
	err := commands.Init()
	if err != nil {
		fmt.Println(err)
	}
}