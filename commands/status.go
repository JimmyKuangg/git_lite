package commands

import (
	"fmt"

	"git_lite/core"
)

func Status() error {
	working, err := core.ScanWorkingDirectory()
	if err != nil {
		return err
	}

	staged, err := core.ScanStagedDirectory()
	if err != nil {
		return err
	}

	if working.Empty() && staged.Empty() {
		fmt.Println("nothing to commit; no changes detected")
		return nil
	}

	printDiff(working, false)
	printDiff(staged, true)

	return nil
}

func printDiff(diff core.Diff, staged bool) {
	if diff.Empty() {
		return
	}

	if staged {
		fmt.Println("Changes to be committed:")
	} else {
		fmt.Println("Changes not staged for commit:")
	}

	if len(diff.Modified) != 0 {
		fmt.Println("  Modified:")
		for _, file := range diff.Modified {
			fmt.Printf("	%s \n", file)
		}
		fmt.Println()
	}

	if len(diff.Untracked) != 0 {
		if staged {
			fmt.Println("  New files:")
		} else {
			fmt.Println("  Untracked Files; Use go run . add [file name or '.'] to add them to the index:")
		}

		for _, file := range diff.Untracked {
			fmt.Printf("	%s \n", file)
		}
		fmt.Println()
	}

	if len(diff.Deleted) != 0 {
		fmt.Println("  Deleted:")
		for _, file := range diff.Deleted {
			fmt.Printf("	%s \n", file)
		}
		fmt.Println()
	}
}
