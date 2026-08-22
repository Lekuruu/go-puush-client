package updater

// Branch identifies the source to download updates from.
type Branch string

const (
	// Regular github releases
	BranchStable Branch = "stable"
	// Pre-release github releases
	BranchPrerelease Branch = "prerelease"
	// Latest github actions builds
	BranchNightly Branch = "nightly"
)
