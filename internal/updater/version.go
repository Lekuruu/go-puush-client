package updater

type Version interface {
	String() string
	CanCompare(other Version) bool
	Compare(other Version) int
	IsNewerThan(other Version) bool
	IsOlderThan(other Version) bool
	IsEqualTo(other Version) bool
}

func NewVersionFromString(versionString string) (Version, error) {
	// For now, we only support semantic versioning
	// We can later try out every versioning scheme we want to support,
	// and return the first one that works
	return NewSemanticVersionFromString(versionString)
}
