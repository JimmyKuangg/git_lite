package commands

import (
	"fmt"
	"git_lite/core"
)

func Status() error {
	statuses, err := core.ScanWorkingDirectory()
	if err != nil {
		return err
	}

	if len(statuses.Deleted) == 0 && len(statuses.Modified) == 0 && len(statuses.Untracked) == 0 {
		fmt.Println("nothing to commit; no changes detected")
		return nil
	}

	if len(statuses.Modified) != 0 {
		fmt.Println("Modified:")
		for _, file := range statuses.Modified {
			fmt.Printf("%s \n", file)
		}
		fmt.Println("")
	}

	if len(statuses.Untracked) != 0 {
		fmt.Println("Untracked:")
		for _, file := range statuses.Untracked {
			fmt.Printf("%s \n", file)
		}
		fmt.Println("")
	}

	if len(statuses.Deleted) != 0 {
		fmt.Println("Deleted:")
		for _, file := range statuses.Deleted {
			fmt.Printf("%s \n", file)
		}
		fmt.Println("")
	}

	return nil
}