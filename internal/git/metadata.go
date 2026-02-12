package git

import "time"

// Metadata represents runtime metadata about repositories and worktrees
type Metadata struct {
	Status      string
	LastCommit  CommitInfo
	AheadBehind AheadBehindInfo
}

// CommitInfo represents commit information
type CommitInfo struct {
	Hash    string
	Message string
	Author  string
	Date    time.Time
}

// AheadBehindInfo represents sync status with remote
type AheadBehindInfo struct {
	Ahead  int
	Behind int
}

// TODO: Implement metadata collection
