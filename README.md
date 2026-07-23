# GitLite

A toy version control system written in Go that recreates a subset of Git's core functionality.

GitLite supports creating repositories, staging files, creating commits, viewing history, and checking out previous commits. The goal wasn't to replace Git. It was to understand how Git actually works by rebuilding it from scratch.

## Disclaimer (O_O)!

Gitlite is intentionally a learning project. It recreates a small subset of Git's functionality with readability and education prioritized over performance or feature completeness.

Not trying to replace Git. Not trying to cause an uprising. Just learning. Don't sue me pls.

## Installation ദ്ദി(˵ •̀ ᴗ - ˵ ) ✧

- Ensure you have [Go installed](https://go.dev)
- Clone the repository and run commands using Go.

```bash
git clone <repo-url>
cd git_lite

go run . init

# Add some files
# Modify them, edit their contents
# Poke around in the .gitlite folder as you do!

go run . add .

go run . commit "Initial commit"

go run . status

go run . log

go run . checkout <commit_hash>
```

## Features 

- Initialize a repository
- Hash and store blob objects
- Build tree objects
- Stage files with an index
- View commit history
- Checkout previous commits

## Why Rebuild Git?

I've used Git for years without really understanding what happens underneath commands like `git add` or `git commit`.

Instead of reading about Git's internals, I wanted to build a simplified version myself. Along the way I discovered that many commands were far more nuanced than I originally thought, especially the staging area. And the whole thing. The whole thing is very nuanced.

## Basic Usage 

```bash
# Initialize a new GitLite repo
go run . init

# Stage files
# Also works with a single local file path
# Can also remove files from staging if you've deleted them
go run . add .

# Commit changes
go run . commit "Initial commit"

# Check status and history
go run . log
go run . status

# Travel back in time
go run . checkout [commit_hash]
```

## Project Structure 

`commands/`

- The actual commands used for the project. Think "add", "status", "log", "commit".

`core/`

- Shared logic used throughout the project, including object storage, hashing, tree construction, commits, and index management.

`.gitlite/`

- Stores the object database used by GitLite. Fun folder to poke around if you're curious how Git stores blobs, trees, and commits.

`docs/`

- Originally intended for development notes. Somewhere along the way those notes turned into blog posts instead. If you're curious about the bugs, misconceptions, and debugging adventures behind GitLite, you can read about them on my portfolio [here](https://jimmy-kuang.com/notes)!

## Limitations

Gitlite intentionally implements only a subset of Git.

Missing features include:

- Branches
- Merge
- Rebase
- Remote Repositories
- Conflict Resolution
- Pack files

## Future Work

If I have the time (or willpower) to improve on this, I may add:

- Branching
- Merging
- Better and cleaner algorithms
- More safety checks

## What I Learned 

This project taught me:

- How git stores blobs, trees, and commits
- How the staging area actually works
- Why git commands often manipulate multiple repository states
- The importance of separating filesystem logic from object storage
- How seemingly simple Git commands hide surprisingly complex behavior

Follow more about what I learned [here in my notes](https://jimmy-kuang.com/notes).

## Deep Dive

If I haven't already advertised it enough yet, you can learn about the madness that was this project (and more) over in my [portfolio](https://jimmy-kuang.com/notes).

There you'll find the bugs, debugging adventures, misconceptions, refactors, and all the moments where I realized I didn't understand Git nearly as well as I thought I did, because that seems to be the story of every step of this project. ദ്ദി╥ ᴗ ╥)
