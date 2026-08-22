package updater

// Branch identifies the source to download updates from.
type Branch string

const (
	// Regular github releases
	BranchStable Branch = "stable"
	// Latest github actions builds
	BranchNightly Branch = "nightly"
)
