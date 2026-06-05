- Trees represent the structure of a directory at a specific point in time
- A tree is a snapshot of a directory, not a live folder
- Trees map names to hashes
- A name can point to either:
  - a blob (file content)
  - another tree (subdirectory)
  ```
    {
      hello.txt : "abc123",
      README.md : "def456",
      docs      : "tree456"
    }
  ```
- Everything within a tree must point to a hash
  - Git always points to the hash of something, acting like a database of immutable objects connected by hashes
