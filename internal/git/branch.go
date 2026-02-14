package git

// Branch represents a Git branch
type Branch struct {
	Name        string
	IsCurrent   bool
	Upstream    string
	Status      string
	DiffAdded   int
	DiffRemoved int
}

// TODO: Implement branch listing
