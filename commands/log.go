package commands

import (
	"fmt"
	"git_lite/core"
	"time"
)

func Log() error {
	_, err := core.EnsureGitliteRepo()
	if err != nil {
		return err
	}

	currentHash, err := core.ReadHEAD()
	if err != nil {
		return err
	}

	for currentHash != "" {
		commit, err := core.ReadCommit(currentHash)
		if err != nil {
			return err
		}

		fmt.Printf("commit %s\n", currentHash)
		fmt.Printf("Author: %s\n", commit.Author)
		fmt.Printf(
    "Date: %s\n",
    commit.Timestamp.Format(time.RFC1123),
		)
		fmt.Printf("\n    %s\n\n", commit.Message)

		currentHash = commit.Parent
	}
	
	return nil
}