package commands

import (
	"git_lite/core"
)

func Commit() error {
	tree, err := core.BuildTree()
	if err != nil {
		return err
	}

	treeHash, err := core.WriteTree(tree)
	if err != nil {
		return err
	}

	commit, err := core.BuildCommit("", treeHash)
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