package commands

import "git_lite/core"

func Checkout(commitHash string) error {
	_, err := core.EnsureGitliteRepo()
	if err != nil {
		return err
	}

	commit, err := core.ReadCommit(commitHash)
	if err != nil {
		return err
	}

	headHash, err := core.ReadHEAD()
	if err != nil {
		return err
	}

	currentCommit, err := core.ReadCommit(headHash)
	if err != nil {
		return err
	}

	currentSnapshot, err := core.FlattenTree(currentCommit.Root)
	if err != nil {
		return err
	}

	treeHash := commit.Root
	snapshot, err := core.FlattenTree(treeHash)
	if err != nil {
		return err
	}
	
	err = core.CleanupRemovedFiles(currentSnapshot, snapshot)
	if err != nil {
		return err
	}

	err = core.ApplySnapshot(snapshot)
	if err != nil {
		return err
	}

	err = core.WriteIndex(snapshot)
	if err != nil {
		return err
	}

	err = core.WriteHEAD(commitHash)
	if err != nil {
		return err
	}

	return nil
}