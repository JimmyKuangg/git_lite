package commands

import "git_lite/core"

func Checkout(commitHash string) error {
	root, err := core.EnsureGitliteRepo()
	if err != nil {
		return err
	}

	commit, err := core.ReadCommit(commitHash)
	if err != nil {
		return err
	}

	working, err := core.BuildWorkingSnapshot(root)
	if err != nil {
		return err
	}

	treeHash := commit.Root
	snapshot, err := core.FlattenTree(treeHash)
	if err != nil {
		return err
	}
	
	err = core.CleanupRemovedFiles(working, snapshot)
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