package main

import (
	"fmt"
	"git_lite/commands"
	"os"
)

func main() {
    args := os.Args
		var err error

		if len(args) < 2 {
			fmt.Println("usage: gitlite <command>")
    	return
		}

    switch args[1] {
    case "add":
      err = commands.Add(args[2])

    case "commit":
      err = commands.Commit(args[2])

    case "log":
      err = commands.Log()
		
		case "init":
			err = commands.Init()
    }

		if err != nil {
			fmt.Println(err)
		}
}