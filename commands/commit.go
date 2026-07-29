package commands

import (
	"fmt"
	"strings"

	"git_lite/core"
)

func Commit(message string) error {
	_, err := core.EnsureGitliteRepo()
	if err != nil {
		return err
	}

	tree, err := core.BuildTree()
	if err != nil {
		return err
	}

	treeHash, err := core.WriteTree(tree)
	if err != nil {
		return err
	}

	// Check HEAD to see if we have a new tree or not before commiting
	headHash, err := core.ReadHEAD()
	if err != nil {
		return err
	}

	// If HEAD is empty, then we have no previous commits
	// No need to compare
	if headHash != "" && !strings.HasPrefix(headHash, "ref: ") {
		prevCommitBytes, err := core.ReadObject(headHash)
		if err != nil {
			return err
		}

		prevCommit, err := core.ParseCommit(prevCommitBytes)
		if err != nil {
			return err
		}

		// Compare the root tree hash of the previous commit with our current tree hash
		// If the two are the same, it means we have made no changes to tracked files
		// Meaning, no need to commit
		if prevCommit.Root == treeHash {
			return fmt.Errorf("nothing to commit; staging area matches HEAD")
		}
	}

	commit, err := core.BuildCommit(message, treeHash)
	if err != nil {
		return err
	}

	commitBytes := commit.Encode()
	commitHash, err := core.WriteObject(commitBytes)
	if err != nil {
		return err
	}

	err = core.WriteHEAD(commitHash)
	if err != nil {
		return err
	}

	return nil
}
