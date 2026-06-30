package main

import (
	"bufio"
	"fmt"
	"git_lite/commands"
	"os"
	"strings"
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
	
	case "status":
		err = commands.Status()

	case "checkout":
		if len(args) < 3 {
      fmt.Println("Usage: go run . checkout <commit-hash>")
      return
    }

		fmt.Println("This process is irreversible and will discard all changes made after this commit. Are you sure?")
		fmt.Println("Continue? (Y/N)")
		
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()

		input := strings.ToUpper(scanner.Text())
		if input != "Y" {
      fmt.Println("Checkout cancelled.")
      return
    }

    err = commands.Checkout(args[2])

	default:
		err = fmt.Errorf("unknown command: '%s'", args[1])
	}
	
	if err != nil {
		fmt.Println(err)
	}
}